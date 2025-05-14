# Remote JavaScript Executor

A tool for watching JavaScript files and executing them in Chrome or Chromium using the Chrome DevTools Protocol.

## Features

- Watch directories for JavaScript file changes and automatically execute them in Chrome
- Execute JavaScript files or inline scripts directly in Chrome
- Start Chrome or Chromium with remote debugging enabled
- Automatic Chrome/Chromium detection on multiple platforms
- Rich logging options for debugging

## Installation

Build the tool using Go:

```bash
cd cmd/apps/remote-js-executor
go build
```

## Usage

### Start Chrome with Remote Debugging

```bash
./remote-js-executor start-chrome --port 9222 [--headless] [--browser-path /path/to/chrome]
```

### Watch JavaScript Files

```bash
./remote-js-executor watch --directory /path/to/watch --pattern "**/*.js" --chrome-url http://localhost:9222
```

### Execute a Single JavaScript File

```bash
./remote-js-executor execute --file /path/to/script.js --chrome-url http://localhost:9222
```

### Execute Inline JavaScript

```bash
./remote-js-executor execute --script "console.log('Hello from Chrome!');" --chrome-url http://localhost:9222
```

## Command Options

### Global Options

- `--log-level` - Set logging level (debug, info, warn, error), default: info

### `start-chrome` Command Options

- `--port` - Port for Chrome remote debugging (default: 9222)
- `--browser-path` - Path to Chrome/Chromium executable (default: auto-detect)
- `--user-data-dir` - User data directory (default: temporary directory)
- `--headless` - Run Chrome in headless mode (default: false)
- `--wait-ms` - Time to wait after starting Chrome in milliseconds (default: 1000)

### `watch` Command Options

- `--directory` - Directory to watch for JavaScript files (required)
- `--pattern` - File pattern to watch (default: "**/*.js")
- `--chrome-url` - Chrome DevTools Protocol URL (default: http://localhost:9222)
- `--navigate-to` - URL to navigate to before executing JavaScript (optional)
- `--throttle-ms` - Throttle execution on rapid changes in milliseconds (default: 500)

### `execute` Command Options

- `--file` - Path to JavaScript file to execute
- `--script` - JavaScript script to execute (alternative to file)
- `--chrome-url` - Chrome DevTools Protocol URL (default: http://localhost:9222)
- `--navigate-to` - URL to navigate to before executing JavaScript (optional)

## Examples

### Full Workflow Example

1. Start Chrome with remote debugging:
   ```bash
   ./remote-js-executor start-chrome --port 9222
   ```

2. Watch a directory for JS changes:
   ```bash
   ./remote-js-executor watch --directory ./scripts --pattern "**/*.js"
   ```

3. Navigate Chrome to a specific page before running scripts:
   ```bash
   ./remote-js-executor watch --directory ./scripts --navigate-to "https://example.com"
   ```

## How It Works

This tool uses the Chrome DevTools Protocol to communicate with Chrome or Chromium. When watching files, it monitors for file changes using Clay's Watcher package and uses chromedp to execute JavaScript code in the browser.