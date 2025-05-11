import { configureStore } from '@reduxjs/toolkit';
import streamReducer from './slices/streamSlice';
import authReducer from './slices/authSlice';
import { streamApi } from '../api/streamApi';
import { stepsApi } from '../api/stepsApi';
import { transcriptApi } from '../api/transcriptApi';
import { githubApi } from '../api/githubApi';
import { authApi } from '../api/authApi';

export const store = configureStore({
  reducer: {
    stream: streamReducer,
    auth: authReducer,
    [streamApi.reducerPath]: streamApi.reducer,
    [stepsApi.reducerPath]: stepsApi.reducer,
    [transcriptApi.reducerPath]: transcriptApi.reducer,
    [githubApi.reducerPath]: githubApi.reducer,
    [authApi.reducerPath]: authApi.reducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      streamApi.middleware,
      stepsApi.middleware,
      transcriptApi.middleware,
      githubApi.middleware,
      authApi.middleware,
    ),
});

// Infer the `RootState` and `AppDispatch` types from the store
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;