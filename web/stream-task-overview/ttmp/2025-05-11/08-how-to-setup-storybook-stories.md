# Guide to Writing Storybook Stories for UI Components

This guide provides a comprehensive overview of how to set up and write Storybook stories for components in this project. It covers mocking Redux state, API interactions using Mock Service Worker (MSW), and best practices to ensure stories are isolated and maintainable.

## 1. Core Concepts and Global Setup

Our Storybook setup is designed to provide a consistent and realistic environment for developing and testing UI components. Key to this is a centralized Redux store and API mocking setup.

### 1.1. Global Redux Store (`ui/.storybook/preview.tsx`)

Instead of each story creating its own Redux store, we have a **global Redux Provider** set up in `ui/.storybook/preview.tsx`. This file configures a single Redux store instance with all the application's reducers and middleware (including RTK Query middleware).

```typescript
// ui/.storybook/preview.tsx (Simplified Excerpt)
import { configureStore } from '@reduxjs/toolkit';
import { Provider } from 'react-redux';
import React from 'react';
// ... import all your reducers and API slices ...
import streamReducer from '../src/store/slices/streamSlice';
import authReducer from '../src/store/slices/authSlice';
import { streamApi } from '../src/api/streamApi';
// ... other api imports

// Export the store for use in stories if direct dispatch is needed (e.g., in play functions)
export const store = configureStore({
  reducer: {
    stream: streamReducer,
    auth: authReducer,
    [streamApi.reducerPath]: streamApi.reducer,
    // ... other reducers and API reducers
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware().concat(
      streamApi.middleware,
      // ... other API middleware
    ),
});

// Global decorator to wrap all stories with the Redux Provider
const withReduxStore = (Story: React.ComponentType) => (
  <Provider store={store}>
    <Story />
  </Provider>
);

const preview: Preview = {
  // ... other parameters ...
  decorators: [withReduxStore, /* other global decorators */],
  loaders: [mswLoader], // MSW loader for API mocking
};

export default preview;
```

**Key takeaway**: You do **not** need to add `<Provider store={...}>` or `configureStore` within your individual `*.stories.tsx` files.

### 1.2. Global API Mocking with MSW

We use Mock Service Worker (MSW) to mock API responses.

*   **MSW Initialization**: `msw-storybook-addon` initializes MSW globally in `ui/.storybook/preview.tsx`.
*   **Global Handlers (`ui/src/mocks/handlers.ts`)**: This file contains default mock API responses for most of our application's endpoints.
    *   **Important**: These handlers **must use absolute URLs** by referencing the `baseUrl` from `ui/src/api/baseApi.ts`. This is crucial for MSW to intercept requests correctly, as our RTK Query setup uses an absolute `baseUrl`.

    ```typescript
    // ui/src/mocks/handlers.ts (Example)
    import { http, HttpResponse } from 'msw';
    import { baseUrl } from '../api/baseApi'; // Crucial import!

    export const handlers = [
      http.get(`${baseUrl}/auth/status`, () => {
        return HttpResponse.json({ isAuthenticated: true, isAdmin: true });
      }),
      http.get(`${baseUrl}/stream`, () => {
        return HttpResponse.json({ title: 'Default Mocked Stream', /* ... */ });
      }),
      // ... other global handlers
    ];
    ```

## 2. Writing Component Stories (`*.stories.tsx`)

Here's a step-by-step guide to creating new stories:

### 2.1. Basic Story Structure

A typical story file looks like this:

```typescript
// src/components/MyComponent.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import MyComponent from './MyComponent';
// Import the global store if you need to dispatch actions in 'play' functions
import { store } from '../../.storybook/preview'; 
// Import actions and initial states from slices
import { someAction, initialState as mySliceInitialState } from '../store/slices/mySlice';
import { resetMySliceState } from '../store/slices/mySlice'; // Assuming you create specific reset actions

const meta: Meta<typeof MyComponent> = {
  title: 'Components/MyComponent',
  component: MyComponent,
  parameters: {
    layout: 'centered', // Or other layout as needed
  },
  tags: ['autodocs'], // Enables automatic documentation generation
};

export default meta;
type Story = StoryObj<typeof MyComponent>;

// Helper function to reset relevant slice(s) to a known state
const resetComponentSpecificStates = () => {
  // Dispatch specific reset actions or set to initial state
  store.dispatch(resetMySliceState(mySliceInitialState)); 
  // store.dispatch({ type: 'MYSLICE_RESET_STATE', payload: mySliceInitialState }); // Alternative if using string actions
};

// Default story
export const Default: Story = {
  play: async () => {
    resetComponentSpecificStates();
    // Optionally dispatch actions to set a specific state for this story
    store.dispatch(someAction({ /* payload */ }));
  },
};

// Story with a specific state
export const ActiveState: Story = {
  play: async () => {
    resetComponentSpecificStates();
    store.dispatch(someAction({ isActive: true }));
  },
};
```

### 2.2. Controlling Redux State

Since the Redux Provider is global, you control the state for your component within a specific story using the `play` function.

1.  **Import the Global Store**:
    `import { store } from '../../.storybook/preview';`
2.  **Import Actions and Initial States**:
    Import any necessary action creators and `initialState` objects from the relevant Redux slices (e.g., `../store/slices/authSlice.ts`, `../store/slices/streamSlice.ts`).
3.  **Reset State**: **Crucially**, to prevent state from one story leaking into another, always reset the relevant parts of the Redux store at the beginning of your `play` function.
    *   You can do this by dispatching specific "reset" actions you define in your slices (preferred for type safety), or by dispatching an action with a predefined initial state payload. We've added `STREAM_RESET_STATE` and `AUTH_RESET_STATE` which are dispatched with `initialState` as payload.

    ```typescript
    // Example from StreamInfoDisplay.stories.tsx
    import { store } from '../../.storybook/preview';
    import { initialState as streamInitialState } from '../store/slices/streamSlice';
    import { initialState as authInitialState, resetAuthState } from '../store/slices/authSlice';

    const resetRelevantStores = () => {
      store.dispatch({ type: 'STREAM_RESET_STATE', payload: streamInitialState }); // Using string type
      store.dispatch(resetAuthState(authInitialState)); // Using action creator
    };

    export const MyStory: Story = {
      play: async () => {
        resetRelevantStores();
        // ... dispatch other actions
      }
    };
    ```
4.  **Dispatch Actions**: After resetting, dispatch actions to set up the specific Redux state your component needs for that particular story.

    ```typescript
    // Example: Setting admin user and edit mode
    import { setCredentials } from '../store/slices/authSlice';
    import { toggleEditMode } from '../store/slices/streamSlice';

    export const AdminEditing: Story = {
      play: async () => {
        resetRelevantStores();
        store.dispatch(setCredentials({ token: 'mock-admin-token', isAdmin: true }));
        store.dispatch(toggleEditMode());
      },
    };
    ```

### 2.3. Mocking API Responses for Specific Stories

Sometimes, a story requires a component to react to a specific API response (e.g., an error state, or a particular data set).

*   Use `parameters.msw.handlers` within the story definition.
*   These story-specific handlers will **override** any global handlers defined in `ui/src/mocks/handlers.ts` for the duration of that story.
*   Remember to use the **absolute `baseUrl`** in these handlers as well.

```typescript
// Example from AuthStatus.stories.tsx
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';

export const UserLoggedIn: Story = {
  play: async () => {
    store.dispatch(resetAuthState(authInitialState));
    // The component will call /api/auth/status, MSW will intercept.
  },
  parameters: {
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return HttpResponse.json({ isAuthenticated: true, isAdmin: false });
        }),
      ],
    },
  },
};

export const AuthError: Story = {
  play: async () => {
    store.dispatch(resetAuthState(authInitialState));
  },
  parameters: {
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return new HttpResponse(null, { status: 500, statusText: 'Server Error' });
        }),
      ],
    },
  },
};
```
If your component fetches data on mount (e.g., via an RTK Query hook like `useGetAuthStatusQuery()`), MSW will intercept the call, and the hook will populate the Redux store with the mocked data (or error). You often don't need to *manually* dispatch `fulfilled` or `rejected` actions for these RTK Query thunks in your stories if the goal is just to mock the API response for the hook. The `play` function is more for setting up pre-conditions or interacting with the component after it renders.

### 2.4. Mocking React Hooks (Rarely Needed Now)

Previously, some components might have directly used custom hooks like `useAuth`. We've been refactoring these to rely on Redux state (e.g., `selectIsAdmin` from `authSlice`). If you encounter a component still using such a hook that needs mocking for a story:

```typescript
// Example of mocking a hook (less common now)
import * as authHooks from '../hooks/useAuth'; // Assuming useAuth is in this module

export const StoryWithMockedHook: Story = {
  decorators: [
    (StoryComponent) => {
      const originalUseAuth = authHooks.useAuth;
      const originalDescriptor = Object.getOwnPropertyDescriptor(authHooks, 'useAuth');
      
      Object.defineProperty(authHooks, 'useAuth', {
        value: () => ({ isAdmin: true, isAuthenticated: true, /* ...other properties */ }),
        configurable: true,
        writable: true,
      });

      // Render the story
      const storyElement = <StoryComponent />;

      // Cleanup: Restore original hook implementation
      if (originalDescriptor) {
        Object.defineProperty(authHooks, 'useAuth', originalDescriptor);
      } else {
        delete (authHooks as any).useAuth; 
      }
      
      return storyElement; // This cleanup timing in a decorator is tricky.
                           // Prefer play functions and Redux state control.
    },
  ],
};
```
**Note**: Direct hook mocking like this in story decorators can be complex to manage for cleanup. Prioritize controlling component behavior via Redux state and API mocking.

## 3. Key Files and Roles Summary

*   **`ui/.storybook/preview.tsx`**: Global Storybook configuration, Redux Provider, MSW addon initialization. Exports the `store`.
*   **`ui/src/mocks/handlers.ts`**: Default global MSW API mock handlers (use absolute URLs).
*   **`ui/src/api/baseApi.ts`**: Defines the `baseUrl` for API calls.
*   **`ui/src/store/index.ts`**: Main application Redux store definition.
*   **`ui/src/store/slices/*.ts`**: Individual Redux slice definitions (actions, reducers, initial state, async thunks). Action creators and `initialState` from here are used in stories.

## 4. Troubleshooting Common Issues

*   **MSW Not Intercepting Requests / 404 Errors for API Calls**:
    *   Ensure your MSW handlers (both global and story-specific) are using the **absolute `baseUrl`** (e.g., `http.get(\`\${baseUrl}/my/endpoint\`, ...)`).
    *   Check for typos in the URL path in your handler vs. the actual API call.
    *   Verify the MSW service worker (`mockServiceWorker.js`) is correctly registered (usually in `public` and initialized via `npx msw init public/`). The `msw-storybook-addon` should handle this, but check browser devtools (Application > Service Workers).
*   **CORS Errors**:
    *   This almost always means MSW is *not* intercepting the request, and your browser is attempting a real cross-origin request to a backend (e.g., `localhost:8080`) that isn't configured for CORS with your Storybook origin (e.g., `localhost:6006`). Double-check handler URLs and MSW setup.
*   **State Leaking Between Stories**:
    *   Always ensure your `play` function for each story (or a shared helper function called by it) resets the relevant Redux store slices to their initial state before applying story-specific state.
*   **TypeScript Errors in Stories**:
    *   Carefully check types for dispatched action payloads.
    *   Ensure you're importing action creators and `initialState` correctly.
    *   For `fetchStreamInfo.pending(null as any, ...)` or similar, ensure the arguments match what the thunk expects (the first argument is `thunkArg`, the second is `meta`). For thunks with `void` `thunkArg`, you might pass `undefined`. `requestId` is the first argument to `pending`, `fulfilled`, and `rejected` actions generated by `createAsyncThunk`.

## 5. Example: Story for `StreamInfoDisplay` (Error State)

```typescript
// ui/src/components/StreamInfoDisplay.stories.tsx (Simplified Error Story)
import type { Meta, StoryObj } from '@storybook/react';
import StreamInfoDisplay from './StreamInfoDisplay';
import { store } from '../../.storybook/preview';
import { fetchStreamInfo, initialState as streamInitialState } from '../store/slices/streamSlice';
import { setCredentials, initialState as authInitialState, resetAuthState } from '../store/slices/authSlice';

const meta: Meta<typeof StreamInfoDisplay> = { /* ... */ };
export default meta;
type Story = StoryObj<typeof StreamInfoDisplay>;

const baseStreamInfo = { /* ... your base stream info object ... */ };

const resetStoreForStreamInfo = () => {
  // Reset stream slice with base info
  store.dispatch({ 
    type: 'STREAM_RESET_STATE', 
    payload: { ...streamInitialState, info: baseStreamInfo } 
  });
  // Reset auth slice
  store.dispatch(resetAuthState(authInitialState));
};

export const ErrorState: Story = {
  play: async () => {
    resetStoreForStreamInfo(); // Reset state first
    // Set user as admin (as an example, if relevant for the component's view)
    store.dispatch(setCredentials({ token: 'mock-token-admin', isAdmin: true }));

    // Dispatch the rejected action for fetchStreamInfo to simulate an API error
    const errorPayload = {
      name: 'StorybookMockError',
      message: 'Failed to load stream information (Storybook mock)',
    };
    // The 'requestId' is a convention for createAsyncThunk actions.
    // The last undefined is for the 'meta' argument of the rejected action.
    store.dispatch(fetchStreamInfo.rejected(errorPayload, 'storybookRequestId', undefined));
  },
};
```

By following these guidelines, you should be able to create robust and maintainable Storybook stories for your components, accurately reflecting different states and API interactions. 