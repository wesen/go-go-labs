# Stream Task Overview Architecture

## Overview

The Stream Task Overview application is a standalone React application designed for live coding streamers. It provides a visually distinctive interface for tracking stream information, tasks, and GitHub integration.

## Tech Stack

- **Frontend**: React 18 with TypeScript
- **State Management**: Redux Toolkit (RTK)
- **Styling**: TailwindCSS
- **Icons**: Lucide React
- **Build Tool**: Vite with Bun

## Project Structure

```
ui/
├── src/
│   ├── components/            # React components
│   │   ├── GithubInfoPanel.tsx  # GitHub integration UI
│   │   ├── StreamInfoDisplay.tsx  # Main component
│   │   ├── TabsNavigation.tsx  # Tab switching UI
│   │   └── TranscriptPanel.tsx  # Stream transcript UI
│   ├── services/              # External service integrations
│   │   └── githubApi.ts      # GitHub API client
│   ├── store/                # Redux state management
│   │   ├── hooks.ts         # TypeScript hooks for Redux
│   │   ├── store.ts         # Redux store configuration
│   │   └── slices/
│   │       └── streamSlice.ts  # Main state slice
│   ├── App.tsx              # Root application component
│   ├── main.tsx             # Entry point
│   └── index.css            # Global styles
├── docs/                    # Documentation
├── public/                  # Static assets
├── index.html               # HTML template
├── package.json             # Dependencies and scripts
└── tailwind.config.js       # TailwindCSS configuration
```

## Component Architecture

### Main Components

1. **StreamInfoDisplay**: The main container component that orchestrates the entire application
2. **GithubInfoPanel**: Handles GitHub repository integration and displays repository information
3. **TabsNavigation**: Provides UI for switching between different views
4. **TranscriptPanel**: Displays the stream transcript with events, tasks, and notes

### Component Hierarchy

```
App
└── StreamInfoDisplay
    ├── Header
    ├── TabsNavigation
    ├── TranscriptPanel (when transcript tab is active)
    └── Main Dashboard (when main tab is active)
        ├── Stream Information Panel
        ├── GithubInfoPanel
        └── Task Management Panel
            ├── Active Task Component
            ├── Completed Tasks List
            └── Upcoming Tasks List
```

## State Management

The application uses Redux Toolkit for state management with a single main slice:

### Stream Slice

The `streamSlice.ts` file contains the primary state for the application, including:

- Stream information (title, description, etc.)
- Task tracking (active, completed, upcoming tasks)
- GitHub integration data
- Stream transcript entries
- UI state (active tab, editing mode, login status)

## Data Flow

1. **User Interactions**: User interacts with the UI components
2. **Actions Dispatched**: Components dispatch Redux actions
3. **State Updates**: Redux reducers update the application state
4. **UI Updates**: Components re-render based on the new state

## Key Features

### Stream Information Management

- Edit stream title, description, language, repository URL
- Track stream duration with automatic updates
- Show viewer count

### Task Tracking System

- Add, complete, and reactivate tasks
- Visualize current progress
- Maintain history of completed tasks

### GitHub Integration

- Connect to GitHub repositories
- Display branch and commit information
- Link to repository and commits

### Stream Content System

- **Notes View**: Record events, tasks, commits, and manually added notes
- **Summary View**: Organized overview with grouped content categories
- **Transcript View**: Raw livestream transcript with speaker attribution
- **LLM Integration**: Support for AI-generated paragraph summaries

## Authentication

The application includes a simulated authentication system that toggles between:

- **Logged In**: Full editing capabilities (streamer view)
- **Logged Out**: Read-only view (viewer view)

See the separate [transcripts-and-notes.md](./transcripts-and-notes.md) document for detailed information about the stream content system.