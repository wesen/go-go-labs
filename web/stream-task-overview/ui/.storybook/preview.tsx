import type { Preview } from "@storybook/react";
import { initialize, mswLoader } from 'msw-storybook-addon';
import '../src/index.css'; // Make sure Tailwind styles are available
import { handlers } from '../src/mocks/handlers';
import { Provider } from 'react-redux';
import React, { useMemo } from 'react';
import { makeStore } from './store';

// This is the msw-storybook-addon specific initialization
// It's separate from our app's MSW setup
initialize({
  onUnhandledRequest: 'bypass',
});

// Create a decorator that creates a fresh Redux store for each story
const withRedux = (Story, ctx) => {
  // Use the preloadedState from parameters.redux if available
  const store = useMemo(
    () => makeStore(ctx.parameters?.redux?.preloadedState),
    [ctx.parameters?.redux?.preloadedState]
  );
  
  return (
    <Provider store={store}>
      <Story {...ctx.args} />
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
  decorators: [withRedux],
};

export default preview;