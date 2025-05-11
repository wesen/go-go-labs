# Stream Notes Feature Documentation

## Overview

The Stream Notes feature provides a chronological record of stream events, including task transitions, GitHub commits, manual notes, and LLM-generated paragraph summaries. It allows streamers to maintain a complete audit trail of stream activities while providing viewers with context about what has happened throughout the session.

## Features

- **Task Tracking**: Automatic recording when tasks are started or completed
- **Commit Integration**: Captures GitHub commits with links to the repository
- **Manual Notes**: Ability to add arbitrary notes to the stream
- **LLM-Generated Paragraphs**: Longer narrative summaries of discussion periods
- **Timestamped Entries**: All entries include precise timestamps
- **Categorized Events**: Color-coded and icon-differentiated event types
- **Sorted Display**: Most recent events shown first
- **Multi-view Interface**: Different ways to view the stream content (Notes, Summary, Transcript)

## Data Schema

### TranscriptEntry Interface

```typescript
interface TranscriptEntry {
  id: string;           // Unique identifier for the entry
  timestamp: string;    // ISO date string of when the event occurred
  type: 'task_started' | 'task_completed' | 'commit' | 'note' | 'paragraph'; // Event type
  content: string;      // Main text content of the entry
  taskName?: string;    // For task events: name of the task
  commitHash?: string;  // For commit events: Git commit hash
  commitUrl?: string;   // For commit events: URL to view the commit
  title?: string;       // For paragraph entries: optional title
  timeRange?: {         // For paragraph entries: time period covered
    start: string;      // Start time (ISO string)
    end: string;        // End time (ISO string)
  };
}
```

### State in Redux

The transcript data is stored in the Redux state as part of the main stream state:

```typescript
interface StreamState {
  // Other state properties...
  transcript: TranscriptEntry[];
  activeTab: 'main' | 'notes' | 'summary' | 'raw';
}
```

## Component Architecture

### Components Hierarchy

```
StreamInfoDisplay
└── TabsNavigation
└── TranscriptPanel
    └── TranscriptEntries
```

### TabsNavigation Component

**File**: `src/components/TabsNavigation.tsx`

- Provides UI for switching between different views (Dashboard, Notes, Summary, Full Text)
- Uses `activeTab` from Redux state
- Dispatches `changeTab` action to update the active tab

### TranscriptPanel Component (Notes View)

**File**: `src/components/TranscriptPanel.tsx`

- Main container for the notes feature
- Displays chronologically sorted entries of all types
- Provides UI for adding manual notes (when logged in)
- Shows appropriate icons and styling for different event types
- Special formatting for LLM-generated paragraph entries
- Handles formatting of timestamps and entry display

### TranscriptSummaryPanel Component (Summary View)

**File**: `src/components/TranscriptSummaryPanel.tsx`

- Provides a structured summary of stream activities
- Organizes entries by type (tasks, commits, notes)
- Displays key statistics and overviews
- Summarizes the stream into meaningful sections

### RawTranscriptPanel Component (Transcript View)

**File**: `src/components/RawTranscriptPanel.tsx`

- Displays the raw transcript of the livestream conversation
- Presents a continuous narrative of the stream
- Converts events into human-readable paragraphs
- Orders events chronologically (oldest first)
- Focuses on readability and flow

## Redux Integration

### Actions

1. **`changeTab`**: Changes the active tab
   ```typescript
   dispatch(changeTab('notes')); // Other options: 'main', 'summary', 'raw' (transcript)
   ```

2. **`addTranscriptNote`**: Adds a manual note
   ```typescript
   dispatch(addTranscriptNote('Started discussing authentication approaches'));
   ```

3. **`addTranscriptParagraph`**: Adds an LLM-generated paragraph summary
   ```typescript
   dispatch(addTranscriptParagraph({
     content: 'Detailed discussion about authentication approaches...',
     title: 'Authentication Discussion',
     timeRange: {
       start: startTime,
       end: endTime
     }
   }));
   ```

### Automatic Event Recording

Transcript entries are automatically generated for certain actions:

1. **Task Started**: When a new active task is set
   ```typescript
   // Inside setNewActiveTopic and completeCurrentStep reducers
   state.transcript.push(generateTranscriptEntry('task_started', 
     `Started working on ${taskName}`, { taskName }));
   ```

2. **Task Completed**: When a task is marked as complete
   ```typescript
   // Inside setNewActiveTopic and completeCurrentStep reducers
   state.transcript.push(generateTranscriptEntry('task_completed', 
     `Completed ${taskName}`, { taskName }));
   ```

3. **Commit Made**: When GitHub commit info is updated
   ```typescript
   // Inside updateGithubCommitInfo reducer
   state.transcript.push(generateTranscriptEntry('commit', 
     commitMessage, { commitHash, commitUrl }));
   ```

### Helper Functions

```typescript
// Helper function to generate a transcript entry
const generateTranscriptEntry = (
  type: TranscriptEntry['type'], 
  content: string, 
  details?: Partial<TranscriptEntry>
): TranscriptEntry => ({
  id: Math.random().toString(36).substr(2, 9),
  timestamp: new Date().toISOString(),
  type,
  content,
  ...details
});
```

## Event Type Styling

| Event Type | Border Color | Icon |
|------------|--------------|------|
| task_started | Blue (#2563eb) | PlayCircle |
| task_completed | Green (#16a34a) | CheckCircle |
| commit | Purple (#9333ea) | GitCommit |
| note | Gray (#6b7280) | Edit3 |
| paragraph | Amber (#d97706) | FileText |

## Usage Examples

### Viewing the Transcript

1. Click on the "Stream Transcript" tab in the navigation
2. The transcript will display all events in reverse chronological order

### Adding a Manual Note

1. Ensure you are in the logged-in state
2. Navigate to the transcript tab
3. Enter text in the note input field
4. Click "Add Note" to add it to the transcript

### Task Tracking

Task events are automatically added when:
- A new active task is set (task_started)
- A task is completed (task_completed)
- A task is reactivated (task_started)

### Commit Integration

Commit events are automatically added when:
- GitHub integration is enabled
- New commit information is received
- The commit entry includes a clickable link to view the commit on GitHub

## Implementation Considerations

### Performance

- The transcript array is immutably updated using Redux
- Sorting happens on render, not in the state
- Rendering optimizations include displaying only the necessary information

### Security

- Note input is sanitized before being added to the transcript
- User must be logged in to add manual notes

### Extensibility

The transcript system can be extended to include additional event types:

1. Add new type to the `TranscriptEntry` type definition
2. Update the `getEntryIcon` function to handle the new type
3. Add styling for the new event type
4. Create actions to add the new event type to the transcript

## Future Enhancements

1. **Filtering**: Allow filtering transcript by event type
2. **Export**: Enable exporting the transcript to Markdown or other formats
3. **Search**: Add search functionality within transcript entries
4. **Rich Content**: Support for rich content like code snippets or images
5. **Categorization**: Allow manual categorization of notes
6. **Linking**: Allow linking entries to specific files or code sections
7. **Time Markers**: Add ability to mark specific points in time for easy reference

## Troubleshooting

### Common Issues

1. **Missing Events**: If events aren't appearing, check that the relevant actions are dispatching correctly
2. **Incorrect Timestamps**: Ensure that timestamps are properly generated as ISO strings
3. **Styling Issues**: Verify that the correct CSS classes are applied to different event types