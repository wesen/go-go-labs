# Frontend-Backend Integration Plan

## Overview

This document outlines the necessary changes to integrate the React frontend (using Redux) with the Go backend API for the Stream Task Overview application.

## Current State Analysis

### Frontend (Redux Store)

- Uses Redux Toolkit for state management
- Contains data models for:
  - Stream information (title, description, etc.)
  - Task steps (completed, active, upcoming)
  - GitHub integration
  - Transcript entries
- Has actions for manipulating state but no API integration

### Backend API

- Provides REST endpoints for:
  - Stream information (GET/PUT)
  - Task steps management (get all, set active, add upcoming, complete, reactivate)
- Uses SQLite for data persistence
- Currently lacks endpoints for transcript/notes and GitHub integration

## Integration Gaps

1. **Data Model Alignment**
   - Backend uses unique IDs for steps, frontend doesn't
   - Data structure differences between frontend state and backend responses

2. **Missing API Connections**
   - No API service layer in frontend to connect Redux actions with backend endpoints
   - No error handling for API failures

3. **Missing Backend Features**
   - Backend lacks transcript/notes management endpoints
   - Backend lacks GitHub integration endpoints

4. **Authentication**
   - Frontend has login state but no authentication mechanism
   - Backend doesn't require authentication yet

## Required Changes

### Frontend Data Model Adaptations

1. **StreamInfo Model Adaptations**
   - Add `id` field to match backend
   - Transform between frontend and backend data structures:
   ```typescript
   // Frontend to Backend transformation
   const toBackendStreamInfo = (frontendInfo: StreamInfo) => ({
     title: frontendInfo.title,
     description: frontendInfo.description,
     startTime: frontendInfo.startTime,
     language: frontendInfo.language,
     githubRepo: frontendInfo.githubRepo,
     viewerCount: frontendInfo.viewerCount
   });

   // Backend to Frontend transformation
   const toFrontendStreamInfo = (backendInfo: any) => ({
     title: backendInfo.title,
     description: backendInfo.description,
     startTime: backendInfo.start_time,  // Note: snake_case to camelCase conversion
     language: backendInfo.language,
     githubRepo: backendInfo.github_repo,
     currentTask: backendInfo.active_step || "",  // Derived from active step
     viewerCount: backendInfo.viewer_count
   });
   ```

2. **Steps Model Adaptations**
   - Update `completedSteps`, `activeStep`, and `upcomingSteps` to handle ID-based operations
   - Implement transformations:
   ```typescript
   // Backend to Frontend transformation for steps
   const toFrontendSteps = (backendSteps: any) => ({
     completedSteps: backendSteps.completed.map((step: any) => step.description),
     activeStep: backendSteps.active ? backendSteps.active.description : "",
     upcomingSteps: backendSteps.upcoming.map((step: any) => step.description),
     // Store the ID mapping for internal use
     _stepIdMapping: {
       ...backendSteps.completed.reduce((acc: any, step: any) => 
         ({ ...acc, [step.description]: step.id }), {}),
       ...backendSteps.upcoming.reduce((acc: any, step: any) => 
         ({ ...acc, [step.description]: step.id }), {}),
       ...(backendSteps.active ? { [backendSteps.active.description]: backendSteps.active.id } : {})
     }
   });
   ```

3. **TranscriptEntry Model Additions**
   - Create backend-compatible structures for all transcript types
   - Handle conversion of timestamp formats if needed

4. **GitHub Integration Model**
   - Add fields required for backend integration
   - Create serialization helpers for GitHub API interactions

### Frontend API Integration

1. **RTK Query Service Definitions**

   **Stream API Service**
   ```typescript
   // streamApi.ts (partial)
   endpoints: (builder) => ({
     getStreamInfo: builder.query<StreamInfo, void>({
       query: () => 'stream',
       transformResponse: (response) => toFrontendStreamInfo(response)
     }),
     updateStreamInfo: builder.mutation<StreamInfo, Partial<StreamInfo>>({
       query: (streamInfo) => ({
         url: 'stream',
         method: 'PUT',
         body: toBackendStreamInfo(streamInfo),
       }),
       transformResponse: (response) => toFrontendStreamInfo(response)
     }),
   })
   ```

   **Steps API Service**
   ```typescript
   // stepsApi.ts (partial)
   endpoints: (builder) => ({
     getAllSteps: builder.query<StepState, void>({
       query: () => 'stream/steps',
       transformResponse: (response) => toFrontendSteps(response)
     }),
     setActiveStep: builder.mutation<StepState, string>({
       query: (step) => ({
         url: 'stream/steps/active',
         method: 'PUT',
         body: { step }
       }),
       transformResponse: (response) => toFrontendSteps(response)
     }),
     addUpcomingStep: builder.mutation<StepState, string>({
       query: (step) => ({
         url: 'stream/steps/upcoming',
         method: 'POST',
         body: { step }
       }),
       transformResponse: (response) => toFrontendSteps(response)
     }),
     completeCurrentStep: builder.mutation<StepState, void>({
       query: () => ({
         url: 'stream/steps/complete',
         method: 'POST'
       }),
       transformResponse: (response) => toFrontendSteps(response)
     }),
     reactivateStep: builder.mutation<StepState, {step: string, source: 'upcoming' | 'completed'}>>({
       query: ({step, source}) => {
         const stepId = getStepIdFromName(step, source);
         return {
           url: 'stream/steps/reactivate',
           method: 'PUT',
           body: { stepId, source }
         };
       },
       transformResponse: (response) => toFrontendSteps(response)
     }),
   })
   ```

   **Transcript API Service** (To be implemented in backend)
   ```typescript
   // transcriptApi.ts (partial)
   endpoints: (builder) => ({
     getTranscript: builder.query<TranscriptEntry[], void>({
       query: () => 'stream/transcript',
     }),
     addNote: builder.mutation<void, string>({
       query: (content) => ({
         url: 'stream/transcript',
         method: 'POST',
         body: { 
           type: 'note',
           content,
           timestamp: new Date().toISOString() 
         }
       })
     }),
     addParagraph: builder.mutation<void, {
       content: string;
       title?: string;
       timeRange?: { start: string; end: string };
     }>({
       query: (data) => ({
         url: 'stream/transcript',
         method: 'POST',
         body: { 
           type: 'paragraph',
           ...data,
           timestamp: new Date().toISOString() 
         }
       })
     }),
     // Additional endpoints for other transcript types
   })
   ```

   **GitHub API Service** (To be implemented in backend)
   ```typescript
   // githubApi.ts (partial)
   endpoints: (builder) => ({
     connectGithub: builder.mutation<void, string>({
       query: (token) => ({
         url: 'github/connect',
         method: 'POST',
         body: { token }
       })
     }),
     getGithubInfo: builder.query<GithubInfo, void>({
       query: () => 'github/info',
     }),
     getCommits: builder.query<any[], void>({
       query: () => 'github/commits',
     }),
   })
   ```

2. **Redux Slice Integration with API**

   The streamSlice.ts needs to be updated to work with these API services:

   ```typescript
   // Async thunks to replace direct state mutations
   export const fetchStreamInfo = createAsyncThunk(
     'stream/fetchStreamInfo',
     async (_, { dispatch }) => {
       const response = await streamApi.endpoints.getStreamInfo.initiate(undefined);
       return response.data;
     }
   );

   export const completeCurrentTask = createAsyncThunk(
     'stream/completeCurrentTask',
     async (_, { dispatch }) => {
       const response = await stepsApi.endpoints.completeCurrentStep.initiate(undefined);
       return response.data;
     }
   );

   // Add more thunks for other operations...
   ```

### Backend API Additions

1. **Transcript API Endpoints**

   **GET /api/stream/transcript**
   - **Purpose**: Retrieve all transcript entries
   - **Response Format**:
   ```json
   [
     {
       "id": "550e8400-e29b-41d4-a716-446655440008",
       "timestamp": "2025-05-11T14:30:45Z",
       "type": "task_started",
       "content": "Started working on project setup",
       "taskName": "Project setup and initialization"
     },
     {
       "id": "550e8400-e29b-41d4-a716-446655440009",
       "timestamp": "2025-05-11T14:35:45Z",
       "type": "note",
       "content": "Created project structure using create-react-app"
     },
     // Other transcript entries...
   ]
   ```

   **POST /api/stream/transcript**
   - **Purpose**: Add a new transcript entry
   - **Request Body**:
   ```json
   {
     "type": "note",
     "content": "Added authentication to API endpoints",
     "timestamp": "2025-05-11T15:30:45Z",
     // Optional fields based on type
     "taskName": "Security implementation",
     "title": "API Security Notes",
     "timeRange": {
       "start": "2025-05-11T15:20:45Z",
       "end": "2025-05-11T15:30:45Z"
     },
     "speaker": "Host",
     "duration": 30
   }
   ```
   - **Response**: The created transcript entry with ID

   **GET /api/stream/transcript/types/:type**
   - **Purpose**: Filter transcript entries by type
   - **Parameters**: `type` - One of: task_started, task_completed, commit, note, paragraph, transcript
   - **Response**: Array of matching transcript entries

2. **GitHub Integration Endpoints**

   **POST /api/github/connect**
   - **Purpose**: Connect to a GitHub repository
   - **Request Body**:
   ```json
   {
     "token": "github_access_token",
     "repoOwner": "yourusername",
     "repoName": "component-library"
   }
   ```
   - **Response**: Success message or error

   **GET /api/github/info**
   - **Purpose**: Get connected GitHub repository information
   - **Response**:
   ```json
   {
     "repoUrl": "https://github.com/yourusername/component-library",
     "isConnected": true,
     "repoOwner": "yourusername",
     "repoName": "component-library",
     "currentBranch": "main",
     "latestCommit": {
       "message": "Added button component",
       "author": "Your Name",
       "hash": "a1b2c3d",
       "date": "2025-05-11T14:30:45Z",
       "url": "https://github.com/yourusername/component-library/commit/a1b2c3d"
     }
   }
   ```

   **GET /api/github/commits**
   - **Purpose**: Get recent commits from the repository
   - **Query Parameters**: `limit` (optional) - Number of commits to return
   - **Response**: Array of commit objects

   **POST /api/github/webhook**
   - **Purpose**: Receive GitHub webhook events
   - **Request Body**: GitHub webhook payload
   - **Response**: Success acknowledgment

### Backend Database Schema Additions

1. **transcript_entries Table**
   ```sql
   CREATE TABLE transcript_entries (
     id TEXT PRIMARY KEY,
     timestamp DATETIME NOT NULL,
     type TEXT NOT NULL,
     content TEXT NOT NULL,
     task_name TEXT,
     commit_hash TEXT,
     commit_url TEXT,
     time_range_start DATETIME,
     time_range_end DATETIME,
     title TEXT,
     speaker TEXT,
     duration INTEGER
   );
   ```

2. **github_integration Table**
   ```sql
   CREATE TABLE github_integration (
     id INTEGER PRIMARY KEY,
     token TEXT NOT NULL,
     repo_owner TEXT NOT NULL,
     repo_name TEXT NOT NULL,
     current_branch TEXT NOT NULL,
     latest_commit_hash TEXT,
     latest_commit_message TEXT,
     latest_commit_author TEXT,
     latest_commit_date DATETIME,
     latest_commit_url TEXT
   );
   ```

3. **Update steps Table**
   ```sql
   -- Ensure steps are stored with IDs
   CREATE TABLE steps (
     id TEXT PRIMARY KEY,
     description TEXT NOT NULL,
     status TEXT NOT NULL, -- 'completed', 'active', or 'upcoming'
     created_at DATETIME NOT NULL,
     completed_at DATETIME
   );
   ```

## Implementation Task List

### Frontend Tasks

- [x] Create api.ts with fetch utilities and base URL configuration
- [x] Create data transformation helpers (toBackendStreamInfo, toFrontendStreamInfo, etc.)
- [x] Implement streamApi.ts using RTK Query for stream endpoints
- [x] Implement stepsApi.ts for task step operations
- [x] Add step ID mapping functionality to handle backend-frontend data structure differences
- [x] Implement transcriptApi.ts for transcript operations
- [x] Implement githubApi.ts for GitHub integration
- [x] Add thunks to streamSlice.ts to use API services instead of direct state mutation
- [x] Update reducers to handle API response structures
- [x] Add loading, error, and success states to the store
- [x] Update components to handle loading/error states
- [x] Create memoized selectors using createSelector for better performance
- [x] Implement ErrorBoundary component for React error handling
- [x] Add mockApiMiddleware for development without a backend

### Backend Tasks

- [x] Create transcript_entries table in SQLite schema
- [x] Create github_integration table in SQLite schema
- [x] Update steps table to use unique IDs and status field
- [x] Implement TranscriptEntry model and repository
- [x] Create transcript API controllers/handlers:
  - [x] GET /api/stream/transcript
  - [x] POST /api/stream/transcript
  - [x] GET /api/stream/transcript/types/:type
- [x] Implement GitHub integration service
- [x] Create GitHub webhook handler and event processing
- [x] Add GitHub API controllers/handlers:
  - [x] POST /api/github/connect
  - [x] GET /api/github/info
  - [x] GET /api/github/commits
  - [x] POST /api/github/webhook
- [x] Modify existing step endpoints to handle ID-based operations
- [x] Add comprehensive error handling and validation
- [x] Create a database seeder utility for development

### Testing Tasks

- [x] Test frontend component rendering with mock data
- [x] Test error handling in frontend components
- [x] Test backend API endpoints
- [x] Test frontend-backend integration for stream info
- [x] Test task step operations (add, complete, reactivate)
- [x] Test transcript operations
- [x] Test GitHub integration
- [x] Test error scenarios and recovery

## API Service Implementation Details

### Example RTK Query Service

```typescript
// api/streamApi.ts
import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react'

export const streamApi = createApi({
  reducerPath: 'streamApi',
  baseQuery: fetchBaseQuery({ baseUrl: 'http://localhost:8080/api' }),
  endpoints: (builder) => ({
    getStreamInfo: builder.query({
      query: () => 'stream',
    }),
    updateStreamInfo: builder.mutation({
      query: (streamInfo) => ({
        url: 'stream',
        method: 'PUT',
        body: streamInfo,
      }),
    }),
    // Add more endpoints here
  }),
})

export const { useGetStreamInfoQuery, useUpdateStreamInfoMutation } = streamApi
```

### Example Redux Integration

```typescript
// store/index.ts
import { configureStore } from '@reduxjs/toolkit'
import streamReducer from './slices/streamSlice'
import { streamApi } from '../api/streamApi'
import { stepsApi } from '../api/stepsApi'
import { transcriptApi } from '../api/transcriptApi'
import { githubApi } from '../api/githubApi'
import { mockApiMiddleware } from '../api/mockApiMiddleware'

export const store = configureStore({
  reducer: {
    stream: streamReducer,
    [streamApi.reducerPath]: streamApi.reducer,
    [stepsApi.reducerPath]: stepsApi.reducer,
    [transcriptApi.reducerPath]: transcriptApi.reducer,
    [githubApi.reducerPath]: githubApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      streamApi.middleware,
      stepsApi.middleware,
      transcriptApi.middleware,
      githubApi.middleware,
      mockApiMiddleware // Add mock middleware to handle CORS issues in development
    ),
})

export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
```

## Timeline

1. **Week 1**: Setup API service layer and basic integration
2. **Week 2**: Implement transcript and GitHub features
3. **Week 3**: Testing and refinement

## Conclusion

By implementing the changes outlined in this document, we will create a fully integrated full-stack application where the React/Redux frontend seamlessly communicates with the Go backend API, providing a cohesive user experience for managing stream tasks and related information.

## Implemented Frontend Features

### API Service Layer
- **BaseApi**: Created common API configuration with fetch utilities
- **StreamApi**: Implemented endpoints for stream information
- **StepsApi**: Implemented endpoints for task step operations
- **TranscriptApi**: Implemented endpoints for transcript operations
- **GithubApi**: Implemented endpoints for GitHub integration

### Data Transformations
- Implemented transformation utilities for converting between frontend and backend data formats
- Added snake_case to camelCase conversion for API responses
- Created step ID mapping functionality to handle backend IDs

### Redux Enhancements
- Added async thunks for API operations with proper error handling
- Updated reducers with extraReducers to handle API responses
- Added loading and error states for all API operations
- Created memoized selectors using createSelector for better performance

### Error Handling & Development
- Implemented ErrorBoundary component for React error handling
- Added mockApiMiddleware to provide mock data during development
- Enhanced components with proper error and loading states
- Used Promise.allSettled for graceful handling of API failures

## Implemented Backend Features

### Database Schema
- Extended SQLite schema with new tables for transcript entries and GitHub integration
- Updated steps table structure to use unique IDs
- Created proper relations between entities

### API Endpoints
- **Transcript API**: Implemented endpoints for managing transcript entries (GET, POST, filter by type)
- **GitHub API**: Implemented endpoints for GitHub repository integration (connect, info, commits, webhook)
- Updated existing step endpoints to handle ID-based operations

### Data Management
- Added comprehensive CRUD operations for all entities
- Implemented proper error handling and validation
- Created transaction support for multi-step operations

### Development Tools
- Created a database seeder utility using `github.com/brianvoe/gofakeit/v6` to populate testing data
- Enhanced logging for debugging and monitoring
- Added CORS support for local development

## Integration Status

The integration of the React/Redux frontend with the Go backend API is now complete. The system can:

1. Fetch and update stream information
2. Manage task steps (add, complete, reactivate)
3. Store and retrieve transcript entries
4. Connect to GitHub repositories and receive webhook events
5. Handle errors gracefully with proper status codes and messages

Both frontend and backend components follow the same data models, making the integration seamless. The system is now ready for production use with a complete full-stack implementation.