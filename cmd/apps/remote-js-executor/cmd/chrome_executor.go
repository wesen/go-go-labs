package cmd

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/rs/zerolog/log"
)

// ChromeExecutor handles communication with Chrome via the DevTools Protocol
type ChromeExecutor struct {
	chromeURL    string
	navigateToURL string
	allocCtx     context.Context
	cancelAlloc  context.CancelFunc
	ctx          context.Context
	cancelCtx    context.CancelFunc
	// mutex ensures safe concurrent access
	mutex        sync.Mutex
}

// NewChromeExecutor creates a new ChromeExecutor instance
func NewChromeExecutor(chromeURL, navigateToURL string) (*ChromeExecutor, error) {
	log.Debug().Str("chromeURL", chromeURL).Str("navigateToURL", navigateToURL).Msg("Creating Chrome executor")

	// Create CDP remote allocator that connects to existing Chrome instance
	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), chromeURL)

	// Create context
	ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(logChromeDebug))

	return &ChromeExecutor{
		chromeURL:    chromeURL,
		navigateToURL: navigateToURL,
		allocCtx:     allocCtx,
		cancelAlloc:  cancelAlloc,
		ctx:          ctx,
		cancelCtx:    cancelCtx,
	}, nil
}

// logChromeDebug is a logger function for Chrome debugging messages
func logChromeDebug(format string, args ...interface{}) {
	log.Debug().Msgf("[ChromeDP] "+format, args...)
}

// ExecuteJavaScript executes the provided JavaScript code in Chrome
func (c *ChromeExecutor) ExecuteJavaScript(ctx context.Context, js string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Create a timeout context
	ctxWithTimeout, cancel := context.WithTimeout(c.ctx, 10*time.Second)
	defer cancel()

	log.Debug().Msg("Preparing to execute JavaScript in Chrome")

	var result string
	actions := []chromedp.Action{}

	// If navigateToURL is set, navigate to that URL first
	if c.navigateToURL != "" {
		log.Debug().Str("url", c.navigateToURL).Msg("Navigating to URL before executing JavaScript")
		actions = append(actions, chromedp.Navigate(c.navigateToURL))
		actions = append(actions, chromedp.Sleep(1*time.Second)) // Give page time to load
	}

	// Add JavaScript execution action
	actions = append(actions, chromedp.Evaluate(js, &result))

	// Run all actions
	if err := chromedp.Run(ctxWithTimeout, actions...); err != nil {
		log.Error().Err(err).Msg("Failed to execute JavaScript")
		return "", fmt.Errorf("failed to execute JavaScript: %w", err)
	}

	log.Debug().Str("result", result).Msg("JavaScript execution successful")
	return result, nil
}

// Close cleans up resources
func (c *ChromeExecutor) Close() {
	log.Debug().Msg("Closing Chrome executor and freeing resources")
	c.cancelCtx()
	c.cancelAlloc()
}