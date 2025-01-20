# Job Reports CLI

Added a new command-line interface for parsing and displaying job reports. This new tool provides both legacy text output and structured output using Glazed, allowing users to easily view and analyze job report data.

- Created `job-reports.go` with `JobReportsCommand` implementation
- Added `main.go` to set up the command-line interface
- Supports parsing multiple report files
- Allows users to choose between summary, job details, or all data
- Provides verbose output option
- Implements both legacy text output and structured Glazed output

# Lua Server Test Files

Added a set of Lua test files to demonstrate various server functionalities.

- Created `lua/` subdirectory with example Lua files
- Added `hello.lua` for basic string responses and query parameter handling
- Added `echo.lua` for echoing request details
- Added `calculator.lua` for simple arithmetic operations
- Added `counter.lua` to demonstrate maintaining state between requests

# Cross-Platform Extension Support

Added Firefox support to the Claude Intercept Extension while maintaining Chrome compatibility.

- Added WebExtension browser API polyfill for cross-browser compatibility
- Updated manifest.json with Firefox-specific settings
- Refactored background.js and popup.js to use browser API
- Updated build system to generate both Chrome and Firefox versions
- Improved popup UI and button handling

# SSE Handler Improvements

Improved the Server-Sent Events (SSE) handler with better goroutine management and modern context handling.

- Replaced deprecated http.CloseNotifier with Request.Context()
- Added errgroup for better goroutine management and error handling
- Improved cleanup logic to prevent channel-related panics
- Streamlined event generation and delivery

# HTMX SSE Integration

Enhanced the Server-Sent Events (SSE) implementation with proper HTMX integration for better real-time updates.

- Updated SSE handler to send both named and unnamed events
- Added HTMX SSE extension support with proper event handling
- Improved HTML template with SSE-specific attributes
- Enhanced styling with modern CSS animations and layout
- Added status updates for connection state

# SSE Start Button

Added manual control over SSE connection initiation for better user experience.

- Added Start Events button to manually initiate SSE connection
- Created new endpoint to handle SSE container initialization
- Improved UI with initial state message and button styling
- Added Clear Events button for better event management
- Enhanced overall user interaction flow

# SSE Detailed Logging

Added comprehensive logging throughout the SSE handler for better debugging and monitoring.

- Added detailed connection lifecycle logging
- Added event generation and sending logs
- Added client-specific logging with remote addresses
- Added event content logging
- Improved error logging with context

# SSE Event Handling Fix

Fixed SSE event handling to properly work with htmx's sse-swap attribute.

- Removed conflicting hx-swap attributes from SSE containers
- Improved event formatting for better compatibility
- Added error handling for event writing
- Ensured proper event delivery to the browser

# SSE Debug Listeners

Added comprehensive client-side debugging for SSE events.

- Added htmx:sseBeforeMessage event listener for swap debugging
- Added htmx:sseMessage event listener for completion tracking
- Added htmx:sseError event listener for error detection
- Added htmx:sseOpen event listener for connection tracking
- Improved SSE container structure for better event handling

# SSE Event Listener Fix

Fixed timing of SSE event listener registration.

- Wrapped event listeners in DOMContentLoaded event
- Ensured event listeners are added after DOM is ready
- Fixed potential race condition in event handling
- Improved reliability of debug logging

# SSE Tutorial Documentation

Added comprehensive tutorial documentation for implementing SSE with HTMX 2.0.

- Created detailed step-by-step guide
- Added code snippets and explanations
- Included best practices and error handling
- Covered both client and server implementation
- Added debugging and troubleshooting tips

# DynamoDB Debug Commands

Added new debug commands for managing DynamoDB job entries in the Textractor application.

- Added `debug dynamo list` command with status and time range filtering
- Added `debug dynamo delete` command for single or batch job deletion
- Implemented filtering by status and age for batch operations
- Added proper error handling and user feedback
- Included pagination support for large result sets

# Fix DynamoDB Reserved Keyword Handling

Fixed an issue where DynamoDB operations would fail when using reserved keywords (like "Error") as attribute names by properly using expression attribute names for all fields in update operations.

- Modified `updateJobStatus` to use expression attribute names for all fields
- Added proper handling of reserved keywords in dynamic field updates
- Ensured consistent handling across both Node.js and Go implementations
- Verified compatibility with existing table schema and indexes
- Improved error handling for DynamoDB update operations
