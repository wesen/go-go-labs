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
	chromeURL     string
	navigateToURL string
	allocCtx      context.Context
	cancelAlloc   context.CancelFunc
	tabCtx        context.Context // persistent tab context
	cancelTab     context.CancelFunc
	mutex         sync.Mutex
	currentURL    string
}

// NewChromeExecutor creates a new ChromeExecutor instance
func NewChromeExecutor(chromeURL, navigateToURL string) (*ChromeExecutor, error) {
	log.Debug().Str("chromeURL", chromeURL).Str("navigateToURL", navigateToURL).Msg("NewChromeExecutor: Creating Chrome executor")

	// Create CDP remote allocator that connects to existing Chrome instance
	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), chromeURL)
	// Check allocCtx immediately
	if allocCtx.Err() != nil {
		log.Warn().Err(allocCtx.Err()).Str("chromeURL", chromeURL).Msg("NewChromeExecutor: allocCtx is already done after NewRemoteAllocator")
		// If allocator context is already bad, we probably should cancel and error out
		cancelAlloc()
		return nil, fmt.Errorf("allocator context failed immediately: %w", allocCtx.Err())
	}

	log.Debug().Msg("NewChromeExecutor: Allocator context created")

	// Create a persistent tab that lives until Close()
	tabCtx, cancelTab := chromedp.NewContext(allocCtx, chromedp.WithLogf(logChromeDebug))
	if tabCtx.Err() != nil {
		cancelAlloc()
		return nil, fmt.Errorf("tab context failed immediately: %w", tabCtx.Err())
	}

	log.Debug().Msg("NewChromeExecutor: Persistent tab context created")

	return &ChromeExecutor{
		chromeURL:     chromeURL,
		navigateToURL: navigateToURL,
		allocCtx:      allocCtx,
		cancelAlloc:   cancelAlloc,
		tabCtx:        tabCtx,
		cancelTab:     cancelTab,
		currentURL:    "",
	}, nil
}

// ExecuteJavaScript executes the provided JavaScript code in Chrome
func (c *ChromeExecutor) ExecuteJavaScript(ctx context.Context, js string) (string, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	log.Debug().Msg("ExecuteJavaScript: Acquired mutex")

	// Create a timeout context for this specific task using the persistent tab context
	log.Debug().Str("js_code", js).Msg("ExecuteJavaScript: Preparing to execute JavaScript")

	var result interface{}
	actions := []chromedp.Action{}

	// If we need to navigate and we're not already at the URL, do it first
	if c.navigateToURL != "" && c.currentURL != c.navigateToURL {
		log.Debug().Str("url", c.navigateToURL).Msg("Navigating to URL before executing JavaScript")
		actions = append(actions, chromedp.Navigate(c.navigateToURL))
		actions = append(actions, chromedp.Sleep(1*time.Second)) // Give page time to load
		c.currentURL = c.navigateToURL                           // Update current URL
	}

	// Add JavaScript execution action
	log.Debug().Msg("ExecuteJavaScript: Adding Evaluate action")
	actions = append(actions, chromedp.Evaluate(js, &result))

	log.Debug().Int("num_actions", len(actions)).Msg("ExecuteJavaScript: About to run actions")

	// Run all actions
	err := chromedp.Run(c.tabCtx, actions...)
	log.Debug().Msg("ExecuteJavaScript: chromedp.Run has completed/returned")

	if err != nil {
		log.Error().Err(err).Str("js_code_on_error", js).Msg("ExecuteJavaScript: Failed during chromedp.Run")
		return "", fmt.Errorf("failed to execute JavaScript: %w", err)
	}

	// Convert the result to a string for logging and return
	var resultStr string
	if result != nil {
		resultStr = fmt.Sprintf("%v", result)
	}

	log.Debug().Str("result", resultStr).Msg("JavaScript execution successful")
	// Print the result to stdout
	fmt.Printf("Result: %s\n", resultStr)
	return resultStr, nil
}

// NavigateToPage navigates Chrome to the specified URL
func (c *ChromeExecutor) NavigateToPage(ctx context.Context, url string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Skip navigation if we're already at the URL
	if c.currentURL == url {
		log.Debug().Str("url", url).Msg("NavigateToPage: Already at the requested URL, skipping navigation")
		return nil
	}

	log.Debug().Str("url", url).Msg("NavigateToPage: Navigating to URL")

	if c.tabCtx.Err() != nil {
		log.Warn().Err(c.tabCtx.Err()).Msg("NavigateToPage: ctxWithTimeout is already done immediately after creation")
		return fmt.Errorf("timeout context creation failed: %w", c.tabCtx.Err())
	}

	// Navigate to the specified URL
	err := chromedp.Run(c.tabCtx,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second), // Give page time to load
	)

	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("NavigateToPage: Failed to navigate to URL")
		return fmt.Errorf("failed to navigate to URL %s: %w", url, err)
	}

	// Update the current URL and the default navigation URL for future calls
	c.currentURL = url
	c.navigateToURL = url

	log.Debug().Str("url", url).Msg("NavigateToPage: Successfully navigated to URL")
	return nil
}

// logChromeDebug is a logger function for Chrome debugging messages
func logChromeDebug(format string, args ...interface{}) {
	log.Debug().Msgf("[ChromeDP] "+format, args...)
}

// Close cleans up resources
func (c *ChromeExecutor) Close() {
	log.Debug().Msg("Closing Chrome executor and freeing resources")
	if c.cancelTab != nil {
		c.cancelTab()
	}
	if c.cancelAlloc != nil {
		c.cancelAlloc()
	}
}
