import type { Meta, StoryObj } from '@storybook/react';
import AuthStatus from './AuthStatus';
import { http, HttpResponse } from 'msw';
import { store } from '../../.storybook/preview';
import { initialState as authInitialState, resetAuthState } from '../store/slices/authSlice';
import { baseUrl } from '../api/baseApi';

const meta: Meta<typeof AuthStatus> = {
  title: 'Components/AuthStatus',
  component: AuthStatus,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof AuthStatus>;

const resetAuthStore = () => {
  store.dispatch(resetAuthState(authInitialState));
};

// Authenticated admin user story
export const AuthenticatedAdmin: Story = {
  play: async () => {
    resetAuthStore();
  },
  parameters: {
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
  play: async () => {
    resetAuthStore();
  },
  parameters: {
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
  play: async () => {
    resetAuthStore();
  },
  parameters: {
    msw: {
      handlers: [
        http.get(`${baseUrl}/auth/status`, () => {
          return HttpResponse.json({ isAuthenticated: false, isAdmin: false });
        })
      ]
    }
  }
};