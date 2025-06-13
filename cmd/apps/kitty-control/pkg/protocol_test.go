package pkg

import (
	"encoding/json"
	"testing"
)

func TestKittyCommand_ToJSON(t *testing.T) {
	cmd := NewKittyCommand("ls")

	data, err := cmd.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	// Verify the JSON structure
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["cmd"] != "ls" {
		t.Errorf("Expected cmd to be 'ls', got %v", result["cmd"])
	}

	if result["version"] == nil {
		t.Error("Expected version to be set")
	}
}

func TestKittyCommand_ToProtocolString(t *testing.T) {
	cmd := NewKittyCommand("ls")

	protocolStr, err := cmd.ToProtocolString()
	if err != nil {
		t.Fatalf("Failed to convert to protocol string: %v", err)
	}

	// Check that it starts and ends with the correct escape sequences
	if len(protocolStr) < 20 {
		t.Error("Protocol string seems too short")
	}

	if protocolStr[:12] != "\x1bP@kitty-cmd" {
		t.Errorf("Protocol string should start with escape sequence")
	}

	if protocolStr[len(protocolStr)-2:] != "\x1b\\" {
		t.Errorf("Protocol string should end with escape sequence")
	}
}

func TestKittyCommand_WithPayload(t *testing.T) {
	payload := &ListPayload{
		AllEnvVars: boolPtr(true),
		Self:       boolPtr(false),
	}

	cmd := NewKittyCommand("ls").WithPayload(payload)

	data, err := cmd.ToJSON()
	if err != nil {
		t.Fatalf("Failed to convert to JSON: %v", err)
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if result["payload"] == nil {
		t.Error("Expected payload to be set")
	}
}

func TestKittyCommand_WithNoResponse(t *testing.T) {
	cmd := NewKittyCommand("send-text").WithNoResponse(true)

	if cmd.NoResponse == nil || *cmd.NoResponse != true {
		t.Error("Expected no_response to be true")
	}
}

// Helper function for tests
func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
