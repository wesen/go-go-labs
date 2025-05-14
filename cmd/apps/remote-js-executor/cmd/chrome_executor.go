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
	// ctx           context.Context // Removed
	// cancelCtx     context.CancelFunc // Removed
	// mutex ensures safe concurrent access
	mutex sync.Mutex
}

// logCancelFunc wraps a context.CancelFunc to log when it is called
func logCancelFunc(name string, cancel context.CancelFunc) context.CancelFunc {
	return func() {
		log.Debug().Str("name", name).Msg("Cancelling context")
		cancel()
	}
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

	// // Create context // Removed this part
	// ctx, cancelCtx := chromedp.NewContext(allocCtx, chromedp.WithLogf(logChromeDebug))
	// // Check ctx immediately
	// if ctx.Err() != nil {
	// 	log.Warn().Err(ctx.Err()).Msg("NewChromeExecutor: ctx is already done after NewContext")
	// }

	log.Debug().Msg("NewChromeExecutor: Allocator context created")

	return &ChromeExecutor{
		chromeURL:     chromeURL,
		navigateToURL: navigateToURL,
		allocCtx:      allocCtx,
		cancelAlloc:   cancelAlloc,
		// ctx:           ctx, // Removed
		// cancelCtx:     cancelCtx, // Removed
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

	log.Debug().Msg("ExecuteJavaScript: Acquired mutex")

	// Create a new task-specific chromedp context from the allocator context
	taskCtx, taskCancel := chromedp.NewContext(c.allocCtx, chromedp.WithLogf(logChromeDebug))
	defer taskCancel()

	if taskCtx.Err() != nil {
		log.Error().Err(taskCtx.Err()).Msg("ExecuteJavaScript: taskCtx from allocator is already done immediately after creation")
		return "", fmt.Errorf("taskCtx creation failed: %w", taskCtx.Err())
	}

	// Create a timeout context for this specific task
	ctxWithTimeout, cancelTimeout := context.WithTimeout(taskCtx, 10*time.Second)
	defer cancelTimeout()

	log.Debug().Str("js_code", js).Msg("ExecuteJavaScript: Preparing to execute JavaScript")

	// // Check parent context status // This referred to c.ctx which is removed
	// if c.ctx.Err() != nil { //
	// 	log.Warn().Err(c.ctx.Err()).Msg("ExecuteJavaScript: Parent context c.ctx is already done before creating timeout context")
	// }
	// Check timeout context status immediately after creation
	if ctxWithTimeout.Err() != nil {
		// This might happen if taskCtx itself was bad, or if 10s is somehow too short (unlikely for just creation)
		log.Warn().Err(ctxWithTimeout.Err()).Msg("ExecuteJavaScript: ctxWithTimeout is already done immediately after creation from taskCtx")
	}

	var result interface{}
	actions := []chromedp.Action{}

	// If navigateToURL is set, navigate to that URL first
	if c.navigateToURL != "" {
		log.Debug().Str("url", c.navigateToURL).Msg("Navigating to URL before executing JavaScript")
		actions = append(actions, chromedp.Navigate(c.navigateToURL))
		actions = append(actions, chromedp.Sleep(1*time.Second)) // Give page time to load
	}

	// Add JavaScript execution action
	log.Debug().Msg("ExecuteJavaScript: Adding Evaluate action")
	actions = append(actions, chromedp.Evaluate(js, &result))

	log.Debug().Int("num_actions", len(actions)).Msg("ExecuteJavaScript: About to run actions")

	// Run all actions
	err := chromedp.Run(ctxWithTimeout, actions...)
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

	log.Debug().Str("url", url).Msg("NavigateToPage: Navigating to URL")

	// Create a new task-specific chromedp context from the allocator context
	taskCtx, taskCancel := chromedp.NewContext(c.allocCtx, chromedp.WithLogf(logChromeDebug))
	defer taskCancel()

	if taskCtx.Err() != nil {
		log.Error().Err(taskCtx.Err()).Msg("NavigateToPage: taskCtx from allocator is already done immediately after creation")
		return fmt.Errorf("taskCtx creation failed: %w", taskCtx.Err())
	}

	// Create a timeout context for this specific task
	ctxWithTimeout, cancelTimeout := context.WithTimeout(taskCtx, 10*time.Second)
	defer cancelTimeout()

	if ctxWithTimeout.Err() != nil {
		log.Warn().Err(ctxWithTimeout.Err()).Msg("NavigateToPage: ctxWithTimeout is already done immediately after creation from taskCtx")
	}

	// Navigate to the specified URL
	err := chromedp.Run(ctxWithTimeout,
		chromedp.Navigate(url),
		chromedp.Sleep(1*time.Second), // Give page time to load
	)

	if err != nil {
		log.Error().Err(err).Str("url", url).Msg("NavigateToPage: Failed to navigate to URL")
		return fmt.Errorf("failed to navigate to URL %s: %w", url, err)
	}

	// Update the default navigation URL for future ExecuteJavaScript calls
	c.navigateToURL = url

	log.Debug().Str("url", url).Msg("NavigateToPage: Successfully navigated to URL")
	return nil
}

// Close cleans up resources
func (c *ChromeExecutor) Close() {
	log.Debug().Msg("Closing Chrome executor and freeing allocator resources")
	// c.cancelCtx() // Removed
	c.cancelAlloc()
}
