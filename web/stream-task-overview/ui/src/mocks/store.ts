import { configureStore } from '@reduxjs/toolkit';
import streamReducer from '../store/slices/streamSlice';
import { stepsApi } from '../api/stepsApi';
import { streamApi } from '../api/streamApi';

// Helper to create a configured store with different states
export const createMockStore = (preloadedState = {}) => {
  return configureStore({
    reducer: {
      stream: streamReducer,
      [stepsApi.reducerPath]: stepsApi.reducer,
      [streamApi.reducerPath]: streamApi.reducer
    },
    middleware: (getDefaultMiddleware) => 
      getDefaultMiddleware()
        .concat(stepsApi.middleware)
        .concat(streamApi.middleware),
    preloadedState
  });
};

// Default mock state for a typical component
export const defaultState = {
  stream: {
    info: {
      title: 'Building a Task Management System',
      description: 'Creating a full-stack application with React and Go',
      language: 'TypeScript/Go',
      githubRepo: 'organization/repo-name',
      startTime: '2025-05-11T18:30:00.000Z',
      viewerCount: 128
    },
    isEditing: false,
    completedSteps: ['Research competitors', 'Create wireframes'],
    activeStep: 'Implement UI components',
    upcomingSteps: ['Write unit tests', 'Deploy to staging'],
    stepIdMapping: {
      'Research competitors': '1',
      'Create wireframes': '2',
      'Implement UI components': '3',
      'Write unit tests': '4',
      'Deploy to staging': '5'
    },
    loading: {
      streamInfo: false,
      steps: false,
      transcript: false,
      github: false
    },
    error: {
      streamInfo: null,
      steps: null,
      transcript: null,
      github: null
    }
  }
};

// Create a store with default mock data
export const mockStore = createMockStore(defaultState);