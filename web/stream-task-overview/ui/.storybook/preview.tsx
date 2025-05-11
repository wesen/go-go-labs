import type { Preview } from "@storybook/react";
import { initialize, mswLoader } from 'msw-storybook-addon';
import '../src/index.css'; // Make sure Tailwind styles are available
import { handlers } from '../src/mocks/handlers';
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import React from 'react';
import streamReducer from '../src/store/slices/streamSlice';
import authReducer from '../src/store/slices/authSlice';
import { streamApi } from '../src/api/streamApi';
import { stepsApi } from '../src/api/stepsApi';
import { transcriptApi } from '../src/api/transcriptApi';
import { githubApi } from '../src/api/githubApi';
import { authApi } from '../src/api/authApi';

// This is the msw-storybook-addon specific initialization
// It's separate from our app's MSW setup
initialize({
  onUnhandledRequest: 'bypass',
});

// Create a store with all reducers and middleware
// Export the store so it can be used in individual stories for dispatching actions
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

// Create a decorator that wraps stories with the Redux provider
const withReduxStore = (Story: React.ComponentType) => {
  return (
    <Provider store={store}>
      <Story />
    </Provider>
  );
};

const preview: Preview = {
  parameters: {
    actions: { argTypesRegex: "^on[A-Z].*" },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    msw: {
      handlers: handlers,
    },
  },
  loaders: [mswLoader],
  decorators: [withReduxStore],
};

export default preview;