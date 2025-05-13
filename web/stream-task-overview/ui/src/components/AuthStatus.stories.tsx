import type { Meta, StoryObj } from '@storybook/react';
import AuthStatus from './AuthStatus';
import { http, HttpResponse } from 'msw';
import { initialState as authInitialState } from '../store/slices/authSlice';
import { baseUrl } from '../api/baseApi';

const meta: Meta<typeof AuthStatus> = {
  title: 'Components/AuthStatus',
  component: AuthStatus,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof AuthStatus>;

// Authenticated admin user story
export const AuthenticatedAdmin: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return HttpResponse.json({ isAuthenticated: true, isAdmin: true });
        })
      ]
    }
  }
};

// Authenticated non-admin user story
export const AuthenticatedUser: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return HttpResponse.json({ isAuthenticated: true, isAdmin: false });
        })
      ]
    }
  }
};

// Not authenticated user story
export const NotAuthenticated: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return HttpResponse.json({ isAuthenticated: false, isAdmin: false });
        })
      ]
    }
  }
};