package pkg

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/rs/zerolog/log"
)

// Client represents a kitty remote control client
type Client struct {
	socketPath string
	timeout    time.Duration
}

// NewClient creates a new kitty client
func NewClient(socketPath string) *Client {
	log.Debug().Str("socket_path", socketPath).Msg("creating new kitty client")
	return &Client{
		socketPath: socketPath,
		timeout:    5 * time.Second,
	}
}

// NewClientFromEnv creates a client using the KITTY_LISTEN_ON environment variable
func NewClientFromEnv() (*Client, error) {
	listenOn := os.Getenv("KITTY_LISTEN_ON")
	log.Debug().Str("kitty_listen_on", listenOn).Msg("checking KITTY_LISTEN_ON environment variable")

	if listenOn == "" {
		log.Error().Msg("KITTY_LISTEN_ON environment variable not set")
		return nil, errors.New("KITTY_LISTEN_ON environment variable not set")
	}

	// Parse the listen address - could be unix:/path/to/socket or tcp:host:port
	if len(listenOn) > 5 && listenOn[:5] == "unix:" {
		socketPath := listenOn[5:]
		log.Debug().Str("socket_path", socketPath).Msg("detected unix socket")
		return NewClient(socketPath), nil
	}

	// For now, only support unix sockets
	log.Error().Str("listen_on", listenOn).Msg("only unix sockets are currently supported")
	return nil, errors.New("only unix sockets are currently supported")
}

// SetTimeout sets the connection timeout
func (c *Client) SetTimeout(timeout time.Duration) {
	log.Debug().Dur("timeout", timeout).Msg("setting client timeout")
	c.timeout = timeout
}

// SendCommand sends a command to kitty and returns the response
func (c *Client) SendCommand(ctx context.Context, cmd *KittyCommand) ([]byte, error) {
	log.Debug().Str("command", cmd.Cmd).Msg("sending command to kitty")

	// Create the protocol string
	protocolString, err := cmd.ToProtocolString()
	if err != nil {
		log.Error().Err(err).Str("command", cmd.Cmd).Msg("failed to convert command to protocol string")
		return nil, fmt.Errorf("failed to convert command to protocol string: %w", err)
	}
	log.Trace().Str("protocol_string", protocolString).Msg("generated protocol string")

	// Connect to the socket
	log.Debug().Str("socket_path", c.socketPath).Dur("timeout", c.timeout).Msg("connecting to kitty socket")
	conn, err := net.DialTimeout("unix", c.socketPath, c.timeout)
	if err != nil {
		log.Error().Err(err).Str("socket_path", c.socketPath).Msg("failed to connect to kitty socket")
		return nil, fmt.Errorf("failed to connect to kitty socket at %s: %w", c.socketPath, err)
	}
	defer func() {
		log.Trace().Msg("closing socket connection")
		conn.Close()
	}()
	log.Debug().Msg("successfully connected to kitty socket")

	// Set deadline for the operation
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(c.timeout)
	}
	log.Trace().Time("deadline", deadline).Msg("setting connection deadline")
	conn.SetDeadline(deadline)

	// Send the command
	log.Trace().Int("payload_size", len(protocolString)).Msg("sending command payload")
	_, err = conn.Write([]byte(protocolString))
	if err != nil {
		log.Error().Err(err).Msg("failed to send command")
		return nil, fmt.Errorf("failed to send command: %w", err)
	}
	log.Debug().Msg("command sent successfully")

	// If no_response is set, don't wait for a response
	if cmd.NoResponse != nil && *cmd.NoResponse {
		log.Debug().Msg("no response expected, returning")
		return nil, nil
	}

	// Read the response
	log.Trace().Msg("reading response from kitty")
	var response []byte
	buffer := make([]byte, 8192)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if n == 0 {
				log.Error().Err(err).Msg("failed to read response")
				return nil, fmt.Errorf("failed to read response: %w", err)
			}
		}

		response = append(response, buffer[:n]...)
		log.Trace().Int("chunk_size", n).Int("total_size", len(response)).Msg("read response chunk")

		// Check if we've reached the end of the kitty protocol message
		if len(response) >= 2 && string(response[len(response)-2:]) == "\x1b\\" {
			log.Debug().Msg("found end of kitty protocol message")
			break
		}

		// If we got less than the buffer size, we've likely reached the end
		if n < len(buffer) {
			log.Debug().Msg("received partial buffer, assuming end of message")
			break
		}
	}

	log.Debug().Int("response_size", len(response)).Msg("received response from kitty")
	log.Trace().Str("response", string(response)).Msg("response content")

	return response, nil
}

// Helper methods for common commands

// List sends a 'ls' command to list windows and tabs
func (c *Client) List(ctx context.Context, payload *ListPayload) ([]byte, error) {
	log.Debug().Msg("listing windows and tabs")
	cmd := NewKittyCommand("ls")
	if payload != nil {
		cmd.WithPayload(payload)
	}
	return c.SendCommand(ctx, cmd)
}

// SendText sends text to a window
func (c *Client) SendText(ctx context.Context, text string, payload *SendTextPayload) error {
	log.Debug().Str("text", text).Msg("sending text to window")
	if payload == nil {
		payload = &SendTextPayload{}
	}
	payload.Data = text

	cmd := NewKittyCommand("send-text").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Str("text", text).Msg("failed to send text")
	}
	return err
}

// FocusWindow focuses a specific window
func (c *Client) FocusWindow(ctx context.Context, match string) error {
	log.Debug().Str("match", match).Msg("focusing window")
	payload := &FocusWindowPayload{Match: &match}
	cmd := NewKittyCommand("focus-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Str("match", match).Msg("failed to focus window")
	}
	return err
}

// Launch launches a new window with the given command
func (c *Client) Launch(ctx context.Context, args []string, payload *LaunchPayload) error {
	log.Debug().Strs("args", args).Msg("launching new window")
	if payload == nil {
		payload = &LaunchPayload{}
	}
	payload.Args = args

	cmd := NewKittyCommand("launch").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Strs("args", args).Msg("failed to launch window")
	}
	return err
}

// CloseWindow closes a window
func (c *Client) CloseWindow(ctx context.Context, match string) error {
	log.Debug().Str("match", match).Msg("closing window")
	payload := &CloseWindowPayload{Match: &match}
	cmd := NewKittyCommand("close-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Str("match", match).Msg("failed to close window")
	}
	return err
}

// CloseWindowWithPayload closes a window with the full payload
func (c *Client) CloseWindowWithPayload(ctx context.Context, payload *CloseWindowPayload) error {
	cmd := NewKittyCommand("close-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// GetText gets text from a window
func (c *Client) GetText(ctx context.Context, payload *GetTextPayload) ([]byte, error) {
	log.Debug().Interface("payload", payload).Msg("getting text from window")
	cmd := NewKittyCommand("get-text")
	if payload != nil {
		cmd.WithPayload(payload)
	}
	response, err := c.SendCommand(ctx, cmd)
	if err != nil {
		log.Error().Err(err).Msg("failed to get text from window")
	}
	return response, err
}

// ResizeWindow resizes a window
func (c *Client) ResizeWindow(ctx context.Context, payload *ResizeWindowPayload) error {
	cmd := NewKittyCommand("resize-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// FocusTab focuses a specific tab
func (c *Client) FocusTab(ctx context.Context, payload *FocusTabPayload) error {
	cmd := NewKittyCommand("focus-tab").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// NewWindow creates a new window
func (c *Client) NewWindow(ctx context.Context, payload *NewWindowPayload) error {
	cmd := NewKittyCommand("new-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// ScrollWindow scrolls a window
func (c *Client) ScrollWindow(ctx context.Context, payload *ScrollWindowPayload) error {
	cmd := NewKittyCommand("scroll-window").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// SetWindowTitle sets the title of a window
func (c *Client) SetWindowTitle(ctx context.Context, payload *SetWindowTitlePayload) error {
	cmd := NewKittyCommand("set-window-title").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// SetTabTitle sets the title of a tab
func (c *Client) SetTabTitle(ctx context.Context, payload *SetTabTitlePayload) error {
	cmd := NewKittyCommand("set-tab-title").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// CreateMarker creates text markers in a window
func (c *Client) CreateMarker(ctx context.Context, payload *CreateMarkerPayload) error {
	cmd := NewKittyCommand("create-marker").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// RemoveMarker removes text markers from a window
func (c *Client) RemoveMarker(ctx context.Context, payload *RemoveMarkerPayload) error {
	cmd := NewKittyCommand("remove-marker").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}

// SetEnv sets environment variables
func (c *Client) SetEnv(ctx context.Context, payload *EnvPayload) error {
	cmd := NewKittyCommand("env").WithPayload(payload).WithNoResponse(true)
	_, err := c.SendCommand(ctx, cmd)
	return err
}
