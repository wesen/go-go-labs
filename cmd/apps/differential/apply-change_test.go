package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestApplyChange(t *testing.T) {
	tests := []struct {
		name        string
		sourceLines []string
		change      Change
		want        []string
		wantErr     error
	}{
		{
			name:        "replace action with valid parameters",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line2", New: "newLine"},
			want:        []string{"line1", "newLine", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Empty source file",
			sourceLines: []string{},
			change:      Change{Action: ActionReplace, Old: "line1", New: "newLine"},
			want:        []string{},
			wantErr:     &ErrCodeBlock{},
		},
		{
			name:        "Replacing whitespace line",
			sourceLines: []string{"\t", "    ", "line3"},
			change:      Change{Action: ActionReplace, Old: "    ", New: "newLine"},
			want:        []string{"\t", "newLine", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Empty target line",
			sourceLines: []string{"", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "", New: "newLine"},
			want:        []string{"newLine", "line2", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Multiple line replacement",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line1\nline2", New: "newLine1\nnewLine2"},
			want:        []string{"newLine1", "newLine2", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Beginning of file replacement",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line1", New: "newLine"},
			want:        []string{"newLine", "line2", "line3"},
			wantErr:     nil,
		},
		{
			name:        "End of file replacement",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line3", New: "newLine"},
			want:        []string{"line1", "line2", "newLine"},
			wantErr:     nil,
		},
		{
			name:        "Non-existent content",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line4", New: "newLine"},
			want:        []string{},
			wantErr:     &ErrCodeBlock{},
		},
		{
			name:        "Mismatch with empty lines",
			sourceLines: []string{"", "", "line3"},
			change:      Change{Action: ActionReplace, Old: "", New: "newLine"},
			want:        []string{"newLine", "", "line3"},
			wantErr:     nil, // or an error if the behavior should be different
		},
		{
			name:        "Exact match requirement",
			sourceLines: []string{"line1", " line2", "line3"}, // Note the whitespace
			change:      Change{Action: ActionReplace, Old: "line2", New: "newLine"},
			want:        []string{},
			wantErr:     &ErrCodeBlock{},
		},
		{
			name:        "Identical old and new content",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line2", New: "line2"}, // no change in the content
			want:        []string{"line1", "line2", "line3"},                       // expected no change in the lines
			wantErr:     nil,                                                       // we expect no error here, as this is a valid operation, though it does nothing
		},
		{
			name:        "Replacing with nothing",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionReplace, Old: "line2", New: ""}, // the content of 'line2' is replaced with an empty string
			want:        []string{"line1", "", "line3"},                       // 'line2' is now empty, effectively deleting the content
			wantErr:     nil,                                                  // no error is expected here as this is a valid operation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyChange(tt.sourceLines, tt.change)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("applyChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// check for nil == []string{}
			if tt.want == nil {
				tt.want = []string{}
			}
			if got == nil {
				got = []string{}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyChange() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestApplyChange_ActionDelete(t *testing.T) {
	tests := []struct {
		name        string
		sourceLines []string
		change      Change
		want        []string
		wantErr     error
	}{
		{
			name:        "Valid deletion",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionDelete, Content: "line2"},
			want:        []string{"line1", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Non-existent content",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionDelete, Content: "line4"},
			want:        nil, // or []string{"line1", "line2", "line3"} if the function does not modify the input on error
			wantErr:     &ErrCodeBlock{},
		},
		{
			name:        "Empty target content",
			sourceLines: []string{"line1", "", "line3"},
			change:      Change{Action: ActionDelete, Content: ""},
			want:        []string{"line1", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Deleting multiple lines",
			sourceLines: []string{"line1", "line2", "line3", "line4"},
			change:      Change{Action: ActionDelete, Content: "line2\nline3"},
			want:        []string{"line1", "line4"},
			wantErr:     nil,
		},
		{
			name:        "Deleting at the beginning",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionDelete, Content: "line1"},
			want:        []string{"line2", "line3"},
			wantErr:     nil,
		},
		{
			name:        "Deleting at the end",
			sourceLines: []string{"line1", "line2", "line3"},
			change:      Change{Action: ActionDelete, Content: "line3"},
			want:        []string{"line1", "line2"},
			wantErr:     nil,
		},
		{
			name:        "Empty source lines",
			sourceLines: []string{},
			change:      Change{Action: ActionDelete, Content: "line1"},
			want:        nil, // or []string{} if the function does not modify the input on error
			wantErr:     &ErrCodeBlock{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := applyChange(tt.sourceLines, tt.change)

			// Error handling: Check if the expected error matches the actual error
			if (err != nil && tt.wantErr != nil && err.Error() != tt.wantErr.Error()) || (err == nil && tt.wantErr != nil) || (err != nil && tt.wantErr == nil) {
				t.Errorf("applyChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// In case of no error, check if the output matches the expected result
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
