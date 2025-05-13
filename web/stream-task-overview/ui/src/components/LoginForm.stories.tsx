import type { Meta, StoryObj } from '@storybook/react';
import LoginForm from './LoginForm';
import { initialState as authInitialState } from '../store/slices/authSlice';
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';

const meta: Meta<typeof LoginForm> = {
  title: 'Components/LoginForm',
  component: LoginForm,
  decorators: [
    (Story) => (
      <div className="max-w-md mx-auto p-4">
        <Story />
      </div>
    ),
  ],
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof LoginForm>;

export const Default: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: {
          ...authInitialState,
          isAuthenticated: false,
          loading: false,
          error: null
        }
      }
    },
    msw: {
      handlers: [
        http.post(`${baseUrl}/auth/login`, async () => {
          return HttpResponse.json({ success: true, token: 'mock-token', isAdmin: true });
        }),
      ],
    },
  }
};

export const WithError: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: {
          ...authInitialState,
          isAuthenticated: false,
          loading: false,
          error: 'Invalid credentials'
        }
      }
    },
    msw: {
      handlers: [
        http.post(`${baseUrl}/auth/login`, async () => {
          return new HttpResponse(
            JSON.stringify({ success: false, message: 'Invalid credentials' }),
            { status: 401 }
          );
        }),
      ],
    },
  }
};

export const Loading: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: {
          ...authInitialState,
          isAuthenticated: false,
          loading: true,
          error: null
        }
      }
    },
  }
};

export const WithCallback: Story = {
  args: {
    onSuccess: () => console.log('Login successful!')
  },
  parameters: {
    redux: {
      preloadedState: {
        auth: {
          ...authInitialState,
          isAuthenticated: false,
          loading: false,
          error: null
        }
      }
    },
  }
};