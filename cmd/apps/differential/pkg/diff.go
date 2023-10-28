package pkg

import (
	"fmt"
	"github.com/go-go-golems/go-go-labs/cmd/apps/differential/kmp"
	"strings"
)

// Change represents a single change in the DSL.
type Change struct {
	Comment          string `json:"comment"`
	Action           Action `json:"action"`
	Old              string `json:"old,omitempty"`
	New              string `json:"new,omitempty"`
	Content          string `json:"content,omitempty"`
	DestinationAbove string `json:"destination_above,omitempty"`
	DestinationBelow string `json:"destination_below,omitempty"`
}

// DSL represents the entire DSL document.
type DSL struct {
	Path    string   `json:"path"`
	Changes []Change `json:"changes"`
}

type Action string

const (
	ActionReplace Action = "replace"
	ActionDelete  Action = "delete"
	ActionMove    Action = "move"
	ActionInsert  Action = "insert"
)

type ErrCodeBlock struct{}

func (e *ErrCodeBlock) Error() string {
	return "specified code block not found in the source"
}

type ErrInvalidChange struct {
	msg string
}

func (e *ErrInvalidChange) Error() string {
	return fmt.Sprintf("invalid change: %s", e.msg)
}

// FindLocation is a function that identifies the position of a specific block of
// code within a given source code. It takes two parameters: sourceLines and
// locationLines and uses KMPSearch to find the matching index.
//
// The function returns two values: the line number (or -1 if not found), and an error
// if the string was not found.
func FindLocation(sourceLines []string, locationLines []string) (int, error) {
	if len(locationLines) == 0 {
		return -1, &ErrCodeBlock{}
	}

	l := kmp.KMPSearch(sourceLines, locationLines)
	if l == -1 {
		return -1, &ErrCodeBlock{}
	}

	return l, nil
}

// ApplyChange applies a specified change to a given set of source lines.
//
// It takes two parameters:
// - sourceLines: A slice of strings representing the source lines to be modified.
// - change: A Change struct detailing the change to be applied.
//
// The function supports four types of actions specified in the Change struct:
// - ActionReplace: Replaces the old content with the new content in the source lines.
// - ActionDelete: Removes the old content from the source lines.
// - ActionMove: Moves the old content to a new location in the source lines.
// - ActionInsert: Inserts new content at a specified location in the source lines.
//
// The function returns a slice of strings representing the modified source lines,
// and an error if the action is unsupported or if there is an issue locating the
// content or destination in the source lines.
func ApplyChange(sourceLines []string, change Change) ([]string, error) {
	switch change.Action {
	case ActionReplace, ActionDelete, ActionMove:
		contentLines := strings.Split(change.Old, "\n")
		if change.Action != ActionReplace {
			contentLines = strings.Split(change.Content, "\n")
		}
		startIdx := kmp.KMPSearch(sourceLines, contentLines)
		if startIdx == -1 {
			return nil, &ErrCodeBlock{}
		}
		endIdx := startIdx + len(contentLines)

		if change.Action == ActionReplace {
			newLines := strings.Split(change.New, "\n")
			sourceLines = append(sourceLines[:startIdx], append(newLines, sourceLines[endIdx:]...)...)
		} else if change.Action == ActionDelete {
			sourceLines = append(sourceLines[:startIdx], sourceLines[endIdx:]...)
		} else if change.Action == ActionMove {
			destination := change.DestinationAbove
			if destination != "" && change.DestinationBelow != "" {
				return nil, &ErrInvalidChange{"Cannot specify both destination_above and destination_below"}
			}
			if destination == "" {
				destination = change.DestinationBelow
			}
			destLines := strings.Split(destination, "\n")
			moveIdx, err := FindLocation(sourceLines, destLines)
			if err != nil {
				return nil, err
			}
			if change.DestinationBelow != "" {
				moveIdx += len(destLines)
			}
			segment := make([]string, endIdx-startIdx)
			copy(segment, sourceLines[startIdx:endIdx])
			sourceLines = append(sourceLines[:startIdx], sourceLines[endIdx:]...)
			if len(sourceLines) < moveIdx {
				sourceLines = append(sourceLines, segment...)
			} else {
				sourceLines = append(sourceLines[:moveIdx], append(segment, sourceLines[moveIdx:]...)...)
			}
		}

	case ActionInsert:
		if len(sourceLines) == 0 {
			sourceLines = append(sourceLines, change.Content)
			break
		}
		contentLines := strings.Split(change.Content, "\n")
		destination := change.DestinationAbove
		if destination != "" && change.DestinationBelow != "" {
			return nil, &ErrInvalidChange{"Cannot specify both destination_above and destination_below"}
		}
		if destination == "" {
			destination = change.DestinationBelow
		}
		destLines := strings.Split(destination, "\n")
		insertIdx, err := FindLocation(sourceLines, destLines)
		if err != nil {
			return nil, err
		}
		if change.DestinationBelow != "" {
			insertIdx += len(destLines)
		}
		sourceLines = append(sourceLines[:insertIdx], append(contentLines, sourceLines[insertIdx:]...)...)

	default:
		return nil, &ErrInvalidChange{"Unsupported action"}
	}

	return sourceLines, nil
}
