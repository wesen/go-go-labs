package pkg

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
)

// KittyCommand represents the basic structure of a kitty remote control command
type KittyCommand struct {
	Cmd           string      `json:"cmd"`
	Version       [3]int      `json:"version"`
	NoResponse    *bool       `json:"no_response,omitempty"`
	KittyWindowID *string     `json:"kitty_window_id,omitempty"`
	Payload       interface{} `json:"payload,omitempty"`
}

// NewKittyCommand creates a new kitty command with standard version
func NewKittyCommand(cmd string) *KittyCommand {
	log.Debug().Str("command", cmd).Msg("creating new kitty command")
	return &KittyCommand{
		Cmd:     cmd,
		Version: [3]int{0, 14, 2}, // Kitty version from the protocol docs
	}
}

// WithNoResponse sets the no_response flag
func (kc *KittyCommand) WithNoResponse(noResponse bool) *KittyCommand {
	log.Trace().Bool("no_response", noResponse).Str("command", kc.Cmd).Msg("setting no_response flag")
	kc.NoResponse = &noResponse
	return kc
}

// WithWindowID sets the kitty_window_id
func (kc *KittyCommand) WithWindowID(windowID string) *KittyCommand {
	log.Trace().Str("window_id", windowID).Str("command", kc.Cmd).Msg("setting window ID")
	kc.KittyWindowID = &windowID
	return kc
}

// WithPayload sets the payload for the command
func (kc *KittyCommand) WithPayload(payload interface{}) *KittyCommand {
	log.Trace().Interface("payload", payload).Str("command", kc.Cmd).Msg("setting command payload")
	kc.Payload = payload
	return kc
}

// ToJSON converts the command to JSON
func (kc *KittyCommand) ToJSON() ([]byte, error) {
	log.Trace().Str("command", kc.Cmd).Msg("converting command to JSON")
	jsonData, err := json.Marshal(kc)
	if err != nil {
		log.Error().Err(err).Str("command", kc.Cmd).Msg("failed to marshal command to JSON")
		return nil, err
	}
	log.Trace().Str("command", kc.Cmd).Int("json_size", len(jsonData)).Msg("successfully converted command to JSON")
	return jsonData, nil
}

// ToProtocolString converts the command to the full kitty protocol string
func (kc *KittyCommand) ToProtocolString() (string, error) {
	log.Trace().Str("command", kc.Cmd).Msg("converting command to protocol string")
	jsonData, err := kc.ToJSON()
	if err != nil {
		return "", err
	}
	protocolString := fmt.Sprintf("\x1bP@kitty-cmd%s\x1b\\", string(jsonData))
	log.Trace().Str("command", kc.Cmd).Int("protocol_size", len(protocolString)).Msg("successfully created protocol string")
	return protocolString, nil
}

// Common payload structures for different commands

// ListPayload for the 'ls' command
type ListPayload struct {
	AllEnvVars *bool   `json:"all_env_vars,omitempty"`
	Match      *string `json:"match,omitempty"`
	MatchTab   *string `json:"match_tab,omitempty"`
	Self       *bool   `json:"self,omitempty"`
}

// SendTextPayload for the 'send-text' command
type SendTextPayload struct {
	Data           string  `json:"data"`
	Match          *string `json:"match,omitempty"`
	MatchTab       *string `json:"match_tab,omitempty"`
	All            *bool   `json:"all,omitempty"`
	ExcludeActive  *bool   `json:"exclude_active,omitempty"`
	SessionID      *string `json:"session_id,omitempty"`
	BracketedPaste *string `json:"bracketed_paste,omitempty"`
}

// FocusWindowPayload for the 'focus-window' command
type FocusWindowPayload struct {
	Match *string `json:"match,omitempty"`
}

// LaunchPayload for the 'launch' command
type LaunchPayload struct {
	Args               []string `json:"args"`
	Match              *string  `json:"match,omitempty"`
	WindowTitle        *string  `json:"window_title,omitempty"`
	Cwd                *string  `json:"cwd,omitempty"`
	Env                []string `json:"env,omitempty"`
	Type               *string  `json:"type,omitempty"`
	KeepFocus          *bool    `json:"keep_focus,omitempty"`
	Location           *string  `json:"location,omitempty"`
	AllowRemoteControl *bool    `json:"allow_remote_control,omitempty"`
	Hold               *bool    `json:"hold,omitempty"`
}

// CloseWindowPayload for the 'close-window' command
type CloseWindowPayload struct {
	Match         *string `json:"match,omitempty"`
	Self          *bool   `json:"self,omitempty"`
	IgnoreNoMatch *bool   `json:"ignore_no_match,omitempty"`
}

// GetTextPayload for the 'get-text' command
type GetTextPayload struct {
	Match          *string `json:"match,omitempty"`
	Extent         *string `json:"extent,omitempty"`
	Ansi           *bool   `json:"ansi,omitempty"`
	Cursor         *bool   `json:"cursor,omitempty"`
	WrapMarkers    *bool   `json:"wrap_markers,omitempty"`
	ClearSelection *bool   `json:"clear_selection,omitempty"`
	Self           *bool   `json:"self,omitempty"`
}

// ResizeWindowPayload for the 'resize-window' command
type ResizeWindowPayload struct {
	Match     *string `json:"match,omitempty"`
	Self      *bool   `json:"self,omitempty"`
	Increment *bool   `json:"increment,omitempty"`
	Axis      *string `json:"axis,omitempty"`
}

// FocusTabPayload for the 'focus-tab' command
type FocusTabPayload struct {
	Match *string `json:"match,omitempty"`
}

// NewWindowPayload for the 'new-window' command
type NewWindowPayload struct {
	Args       []string `json:"args,omitempty"`
	Match      *string  `json:"match,omitempty"`
	Title      *string  `json:"title,omitempty"`
	Cwd        *string  `json:"cwd,omitempty"`
	KeepFocus  *bool    `json:"keep_focus,omitempty"`
	WindowType *string  `json:"window_type,omitempty"`
	NewTab     *bool    `json:"new_tab,omitempty"`
	TabTitle   *string  `json:"tab_title,omitempty"`
}

// ScrollWindowPayload for the 'scroll-window' command
type ScrollWindowPayload struct {
	Amount []interface{} `json:"amount,omitempty"`
	Match  *string       `json:"match,omitempty"`
}

// SetWindowTitlePayload for the 'set-window-title' command
type SetWindowTitlePayload struct {
	Title     *string `json:"title,omitempty"`
	Match     *string `json:"match,omitempty"`
	Temporary *bool   `json:"temporary,omitempty"`
}

// SetTabTitlePayload for the 'set-tab-title' command
type SetTabTitlePayload struct {
	Title string  `json:"title"`
	Match *string `json:"match,omitempty"`
}

// CreateMarkerPayload for the 'create-marker' command
type CreateMarkerPayload struct {
	Match      *string  `json:"match,omitempty"`
	Self       *bool    `json:"self,omitempty"`
	MarkerSpec []string `json:"marker_spec,omitempty"`
}

// RemoveMarkerPayload for the 'remove-marker' command
type RemoveMarkerPayload struct {
	Match *string `json:"match,omitempty"`
	Self  *bool   `json:"self,omitempty"`
}

// EnvPayload for the 'env' command
type EnvPayload struct {
	Env map[string]string `json:"env"`
}
