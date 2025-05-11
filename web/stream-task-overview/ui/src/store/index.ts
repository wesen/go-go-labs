import { configureStore } from '@reduxjs/toolkit';
import streamReducer from './slices/streamSlice';
import { streamApi } from '../api/streamApi';
import { stepsApi } from '../api/stepsApi';
import { transcriptApi } from '../api/transcriptApi';
import { githubApi } from '../api/githubApi';
import { mockApiMiddleware } from '../api/mockApiMiddleware';

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
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;