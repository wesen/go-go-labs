import { configureStore } from '@reduxjs/toolkit';
import streamReducer from '../src/store/slices/streamSlice';
import authReducer from '../src/store/slices/authSlice';
import { streamApi } from '../src/api/streamApi';
import { authApi } from '../src/api/authApi';
import { transcriptApi } from '../src/api/transcriptApi';
import { stepsApi } from '../src/api/stepsApi';
import { githubApi } from '../src/api/githubApi';

export const makeStore = (preloadedState = {}) => {
  // Ensure we have default values for important properties
  const initialStore = configureStore({
    reducer: {
      stream: streamReducer,
      auth: authReducer,
      [streamApi.reducerPath]: streamApi.reducer,
      [authApi.reducerPath]: authApi.reducer,
      [transcriptApi.reducerPath]: transcriptApi.reducer,
      [stepsApi.reducerPath]: stepsApi.reducer,
      [githubApi.reducerPath]: githubApi.reducer,
    },
    middleware: (getDefaultMiddleware) =>
      getDefaultMiddleware().concat(
        streamApi.middleware,
        authApi.middleware,
        transcriptApi.middleware,
        stepsApi.middleware,
        githubApi.middleware
      ),
    preloadedState,
  });
  
  // Log the initial state for debugging
  console.log('Store created with initial state:', initialStore.getState());
  
  return initialStore;
};