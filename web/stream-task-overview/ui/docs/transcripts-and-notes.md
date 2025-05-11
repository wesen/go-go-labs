# Stream Transcripts and Notes Documentation

## Overview

The Stream Task Overview application includes comprehensive features for tracking stream content through three main views:

1. **Notes** - Chronological log of events, activities, observations, and manually added notes
2. **Summary** - Organized, categorized summary of stream content with key highlights
3. **Transcript** - Raw transcript of the livestream conversation with speaker attribution

## Data Schema

### Core Entry Type

All stream content is stored using the `TranscriptEntry` interface:

```typescript
interface TranscriptEntry {
  id: string;           // Unique identifier for the entry
  timestamp: string;    // ISO date string of when the event occurred
  type: string;         // Type classifier for the entry
  content: string;      // Main text content of the entry
  
  // Optional fields based on entry type
  taskName?: string;    // For task events: name of the task
  commitHash?: string;  // For commit events: Git commit hash
  commitUrl?: string;   // For commit events: URL to view the commit
  
  // For LLM-generated paragraphs
  title?: string;       // Optional paragraph title 
  timeRange?: {         // Time period covered
    start: string;      // Start time (ISO string)
    end: string;        // End time (ISO string)
  };
  
  // For transcript entries
  speaker?: string;     // Name of the speaker
  duration?: number;    // Duration in seconds
}
```

### Entry Types

The `type` field classifies entries into one of these categories:

| Type | Description | Special Fields | Visual |  
|------|-------------|----------------|--------|  
| `task_started` | Marks when a task begins | `taskName` | Blue border/icon |
| `task_completed` | Records task completion | `taskName` | Green border/icon |
| `commit` | Records a GitHub commit | `commitHash`, `commitUrl` | Purple border/icon |
| `note` | Manual note added during stream | - | Gray border/icon |
| `paragraph` | LLM-generated summary of discussion | `title`, `timeRange` | Amber border/icon |
| `transcript` | Raw livestream transcript snippet | `speaker`, `duration` | Default formatting |

## Components and Views

### 1. Notes View (`TranscriptPanel.tsx`)

Displays a chronological log of all events and activities, with the most recent on top. Focuses on presenting individual events with their timestamps and metadata.

- Shows task transitions, commits, manual notes, and paragraph summaries
- Visually differentiates entry types with colors and icons
- Includes a form for adding manual notes (when logged in)

### 2. Summary View (`TranscriptSummaryPanel.tsx`)

Organizes entries by type to provide a structured overview of the stream. Groups content into meaningful sections for easier comprehension.

- **Stream Overview** - Quick summary of stream activities
- **Task Progression** - Current status and completed tasks
- **Code Changes** - Summary of GitHub commits
- **Stream Notes** - Highlights from manually added notes

### 3. Transcript View (`RawTranscriptPanel.tsx`)

Provides the verbatim transcript of the livestream conversation, organized chronologically. Each entry includes:

- Speaker name
- Timestamp 
- Duration (when available)
- Spoken content

## Adding Content

### 1. Automatic Entries

The system automatically adds entries when:

- A task is started or completed
- A GitHub commit is registered

### 2. Manual Notes

Users in the logged-in state can add notes through the Notes view. These appear instantly in the stream chronology.

### 3. LLM-Generated Paragraphs

Longer narrative summaries can be added through the API. These typically summarize discussion periods and are created by processing the raw transcript with an LLM.

### 4. Raw Transcript Entries

These are added through the backend API by integrating with streaming platforms or audio transcription services.

## Redux Actions

```typescript
// Change the active view
dispatch(changeTab('notes')); // Options: 'main', 'notes', 'summary', 'raw'

// Add a manual note
dispatch(addTranscriptNote('We decided to use CSS modules for styling components'));

// Add an LLM-generated paragraph summary
dispatch(addTranscriptParagraph({
  content: 'Detailed discussion about styling approaches...',
  title: 'Component Styling Discussion',
  timeRange: {
    start: startTimestamp,
    end: endTimestamp
  }
}));

// Add a raw transcript entry
dispatch(addTranscriptEntry({
  content: 'I think we should consider server components for this feature.',
  timestamp: new Date().toISOString(),
  speaker: 'Guest',
  duration: 8 // seconds
}));
```

## Integration with External Services

### Twitch/YouTube Integration

The transcript view can be populated with entries from streaming platforms by:

1. Using platform APIs to retrieve chat messages and events
2. Converting audio to text with a transcription service
3. Adding each segment using the `addTranscriptEntry` action

### GitHub Integration

Commit information automatically appears in the Notes and Summary views. The system:

1. Connects to GitHub using a personal access token
2. Monitors for new commits
3. Adds commit entries with links to the repository

## Future Enhancements

- **Filtering**: Allow filtering entries by type in all views
- **Search**: Add search functionality across all content
- **Export**: Enable exporting the transcript, notes, or summary
- **Rich Media**: Support for embedded screenshots or code snippets
- **Speaker Analytics**: Visualizations of speaker participation
- **Timeline View**: Interactive timeline of stream activities