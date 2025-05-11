# Guide to Writing Storybook Stories for UI Components

This guide provides a comprehensive overview of how to set up and write Storybook stories for components in this project. It covers mocking Redux state, API interactions using Mock Service Worker (MSW), and best practices to ensure stories are isolated and maintainable.

## 1. Core Concepts and Setup

Our Storybook setup is designed to provide a consistent and realistic environment for developing and testing UI components. 

### 1.1. Redux Store Factory (Recommended Approach)

To prevent state from bleeding between stories, we use a **store factory** approach that creates a fresh Redux store for each story.

```typescript
// ui/.storybook/store.ts
import { configureStore } from '@reduxjs/toolkit';
import rootReducer from '../src/store/rootReducer';
import { streamApi } from '../src/api/streamApi';

export const makeStore = (preloadedState = {}) =>
  configureStore({
    reducer: {
      ...rootReducer,
      [streamApi.reducerPath]: streamApi.reducer,
    },
    middleware: gDM => gDM().concat(streamApi.middleware),
    preloadedState,
  });
```

Our Redux Provider is set up in `ui/.storybook/preview.tsx` using this factory:

```typescript
// ui/.storybook/preview.tsx
import React, { useMemo } from 'react';
import { Provider } from 'react-redux';
import { Preview } from '@storybook/react';
import { initialize, mswLoader } from 'msw-storybook-addon';
import { makeStore } from './store';

initialize();      // MSW

const withRedux = (Story, ctx) => {
  // optional `parameters.redux.preloadedState`
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
  decorators: [withRedux],
  loaders: [mswLoader],
};

export default preview;
```

**Key benefit**: Each story gets its own isolated Redux store, preventing state bleeding between stories in the Docs page or when multiple stories are rendered side-by-side.

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

### 2.1. Declarative Story Structure

With the store factory approach, you can define story state declaratively using parameters:

```typescript
// src/components/MyComponent.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import MyComponent from './MyComponent';

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

// Default story
export const Default: Story = {
  args: { /* component props */ },
  parameters: {
    redux: {
      preloadedState: {
        auth: { isAuthenticated: false },
        stream: { /* initial state */ },
      },
    },
  },
};

// Story with a specific state
export const ActiveState: Story = {
  args: { /* component props */ },
  parameters: {
    redux: {
      preloadedState: {
        auth: { isAuthenticated: true, isAdmin: true },
        stream: { editMode: true, /* other state */ },
      },
    },
  },
};
```

**No need for complex reset logic**: Since each story gets a fresh store instance, you simply declare the exact state you want via `parameters.redux.preloadedState`.

### 2.2. Controlling Redux State

With this approach, you define the Redux state declaratively through `parameters`:

1. **Declare Initial State**: 
   Use `parameters.redux.preloadedState` to set the exact initial state for each story.

   ```typescript
   // Example from StreamInfoDisplay.stories.tsx
   export const AdminView: Story = {
     parameters: {
       redux: {
         preloadedState: {
           auth: { isAuthenticated: true, isAdmin: true },
           stream: { info: { title: 'Test Stream', /* other props */ } },
         },
       },
     },
   };
   ```

2. **For Interactive Stories**: 
   If you need to interact with the component after initial render, you can still use `play` functions. Each story's `play` function will operate on its own isolated store:

   ```typescript
   // Example: Interacting with component after render
   export const ToggleEditMode: Story = {
     parameters: {
       redux: {
         preloadedState: {
           auth: { isAuthenticated: true, isAdmin: true },
           stream: { editMode: false },
         },
       },
     },
     play: async ({ canvasElement, step }) => {
       const canvas = within(canvasElement);
       
       // Find and click the edit button
       await step('Click edit button', async () => {
         const editButton = canvas.getByRole('button', { name: /edit/i });
         await userEvent.click(editButton);
       });
       
       // Verify the component is in edit mode (store is updated)
       await step('Verify edit mode', async () => {
         // Component should now show edit controls
         expect(canvas.getByRole('button', { name: /save/i })).toBeInTheDocument();
       });
     },
   };
   ```

### 2.3. Mocking API Responses for Specific Stories

MSW setup works just as before - each story can override global handlers:

```typescript
// Example from AuthStatus.stories.tsx
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';

export const UserLoggedIn: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { isAuthenticated: false }, // Initial state before API call
      },
    },
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
  parameters: {
    redux: {
      preloadedState: {
        auth: { isAuthenticated: false },
      },
    },
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

If your component fetches data on mount (e.g., via an RTK Query hook like `useGetAuthStatusQuery()`), MSW will intercept the call, and the hook will populate each story's isolated Redux store with the mocked data (or error).

### 2.4. Mocking React Hooks (Rarely Needed Now)

The approach for mocking hooks remains the same as before. However, with components now relying more on Redux, this is rarely needed.

## 3. Key Files and Roles Summary

*   **`ui/.storybook/store.ts`**: Redux store factory that creates fresh store instances.
*   **`ui/.storybook/preview.tsx`**: Global Storybook configuration with per-story Redux Provider setup.
*   **`ui/src/mocks/handlers.ts`**: Default global MSW API mock handlers (use absolute URLs).
*   **`ui/src/api/baseApi.ts`**: Defines the `baseUrl` for API calls.
*   **`ui/src/store/rootReducer.ts`**: Main application Redux reducers.
*   **`ui/src/store/slices/*.ts`**: Individual Redux slice definitions (actions, reducers, initial state, async thunks).

## 4. Why Per-Story Stores Matter

### The Problem with a Single Global Store

Storybook shows more than one story at the same time:

* The **Canvas** tab
* The **Docs** page, which often renders **all** variants (Primary, Secondary, etc.) side-by-side inside one iframe

With a single global store instance:

* State that a `play()` function mutates for one variant is immediately visible in the other variants rendered right next to it.
* When you embed several components (for example with MDX's `<Canvas>` blocks) they cannot display independent states.

Storybook's own docs point out that code in `preview.tsx` "runs for **every story**" and therefore affects every canvas that is rendered together.

### Advantages of Per-Story Stores

| Old approach (one global store)                      | New per-story store                                |
| ---------------------------------------------------- | -------------------------------------------------- |
| State bleeds between stories in Docs                 | Every story is sandboxed                           |
| Complex `play()` functions to reset slices           | Pure declarative `parameters.redux.preloadedState` |
| Can't open two different auth scenarios side-by-side | Safe to show *any* mix of scenarios together       |
| Hard to reason about RTK Query cache                 | Fresh RTK Query cache per story                    |

## 5. Special Case: Shared Store (When Needed)

If you absolutely need a shared store for special interactive stories, create a separate story file (e.g., `ConnectedPlayground.stories.tsx`) that deliberately re-uses the same store instance:

```typescript
// ConnectedPlayground.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import { Provider } from 'react-redux';
import ConnectedComponents from './ConnectedComponents';
import { configureStore } from '@reduxjs/toolkit';
// ... other imports

// Create a single shared store instance for this file only
const sharedStore = configureStore({
  reducer: { /* ... */ },
  middleware: /* ... */,
});

const meta: Meta<typeof ConnectedComponents> = {
  title: 'Playground/ConnectedComponents',
  component: ConnectedComponents,
  decorators: [
    (Story) => (
      <Provider store={sharedStore}>
        <Story />
      </Provider>
    ),
  ],
  // Prevent these stories from appearing in the docs
  tags: ['!docs'],
};

export default meta;
type Story = StoryObj<typeof ConnectedComponents>;

// Stories in this file will share the same store instance
export const FirstScenario: Story = { /* ... */ };
export const SecondScenario: Story = { /* ... */ };
```

This approach gives you the best of both worlds - isolated stores for documentation and a shared store when you specifically need interaction between components.

## 6. Example: Declarative Story for `StreamInfoDisplay`

```typescript
// ui/src/components/StreamInfoDisplay.stories.tsx
import type { Meta, StoryObj } from '@storybook/react';
import StreamInfoDisplay from './StreamInfoDisplay';
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';

const meta: Meta<typeof StreamInfoDisplay> = {
  title: 'Components/StreamInfoDisplay',
  component: StreamInfoDisplay,
  tags: ['autodocs'],
};
export default meta;
type Story = StoryObj<typeof StreamInfoDisplay>;

export const Loading: Story = {
  parameters: {
    redux: {
      preloadedState: {
        stream: { loading: true },
      },
    },
  },
};

export const ErrorState: Story = {
  parameters: {
    redux: {
      preloadedState: {
        stream: { 
          loading: false,
          error: 'Failed to load stream information (Storybook mock)',
        },
        auth: { isAuthenticated: true, isAdmin: true },
      },
    },
  },
};

export const AdminView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        stream: { 
          loading: false, 
          info: { 
            title: 'Live Coding Session',
            description: 'Building a React application with Redux',
            startTime: '2025-05-11T15:00:00Z',
          },
        },
        auth: { isAuthenticated: true, isAdmin: true },
      },
    },
  },
};

export const UserView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        stream: { 
          loading: false, 
          info: { 
            title: 'Live Coding Session',
            description: 'Building a React application with Redux',
            startTime: '2025-05-11T15:00:00Z',
          },
        },
        auth: { isAuthenticated: true, isAdmin: false },
      },
    },
  },
};
```

By following these updated guidelines, you'll create robust, isolated, and maintainable Storybook stories for your components that accurately reflect their various states without interference between stories. 