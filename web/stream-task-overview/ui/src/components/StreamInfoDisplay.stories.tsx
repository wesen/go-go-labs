import type { Meta, StoryObj } from '@storybook/react';
import StreamInfoDisplay from './StreamInfoDisplay';
import { initialState as streamInitialState } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';

const meta: Meta<typeof StreamInfoDisplay> = {
  title: 'Components/StreamInfoDisplay',
  component: StreamInfoDisplay,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof StreamInfoDisplay>;

const baseStreamInfo = {
  title: 'Building a Task Management System',
  description: 'Creating a full-stack application with React and Go',
  language: 'TypeScript/Go',
  githubRepo: 'organization/repo-name',
  startTime: '2025-05-11T18:30:00.000Z',
  viewerCount: 128,
  currentTask: 'Default Task'
};

export const Default: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: { 
          ...streamInitialState, 
          info: baseStreamInfo,
          loading: { ...streamInitialState.loading, streamInfo: false },
          error: { ...streamInitialState.error, streamInfo: null }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream`, () => {
          return HttpResponse.json(baseStreamInfo);
        }),
      ],
    },
  },
};

export const UserView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token-user' },
        stream: { 
          ...streamInitialState, 
          info: baseStreamInfo,
          loading: { ...streamInitialState.loading, streamInfo: false },
          error: { ...streamInitialState.error, streamInfo: null }
        },
      },
    },
  },
};

export const EditMode: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: { 
          ...streamInitialState, 
          info: baseStreamInfo,
          isEditing: true,
          loading: { ...streamInitialState.loading, streamInfo: false },
          error: { ...streamInitialState.error, streamInfo: null }
        },
      },
    },
  },
};

export const Loading: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: { 
          ...streamInitialState, 
          loading: { ...streamInitialState.loading, streamInfo: true },
          error: { ...streamInitialState.error, streamInfo: null }
        },
      },
    },
  },
};

export const Error: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: { 
          ...streamInitialState, 
          loading: { ...streamInitialState.loading, streamInfo: false },
          error: { ...streamInitialState.error, streamInfo: 'Failed to load stream information (Storybook mock)' }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream`, () => {
          return new HttpResponse(null, { status: 500, statusText: 'Server Error' });
        }),
      ],
    },
  },
};