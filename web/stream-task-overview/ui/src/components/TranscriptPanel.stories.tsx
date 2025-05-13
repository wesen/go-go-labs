import type { Meta, StoryObj } from '@storybook/react';
import TranscriptPanel from './TranscriptPanel';
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';
import {
  initialState as streamInitialState,
  TranscriptEntry
} from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof TranscriptPanel> = {
  title: 'Components/TranscriptPanel',
  component: TranscriptPanel,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof TranscriptPanel>;

// Sample transcript entries for different states
const sampleTranscript: TranscriptEntry[] = [
  {
    id: '1',
    timestamp: new Date(Date.now() - 45 * 60000).toISOString(),
    type: 'task_started',
    content: 'Started working on project setup',
    taskName: 'Project setup'
  },
  {
    id: '2',
    timestamp: new Date(Date.now() - 40 * 60000).toISOString(),
    type: 'note',
    content: 'Created initial project structure'
  },
  {
    id: '3',
    timestamp: new Date(Date.now() - 30 * 60000).toISOString(),
    type: 'commit',
    content: 'Initial project setup',
    commitHash: 'a1b2c3d',
    commitUrl: 'https://github.com/org/repo/commit/a1b2c3d'
  },
  {
    id: '4',
    timestamp: new Date(Date.now() - 20 * 60000).toISOString(),
    type: 'paragraph',
    title: 'Discussion Summary',
    content: 'We discussed the project structure and decided on using React with TypeScript.',
    timeRange: {
      start: new Date(Date.now() - 25 * 60000).toISOString(),
      end: new Date(Date.now() - 15 * 60000).toISOString()
    }
  },
  {
    id: '5',
    timestamp: new Date(Date.now() - 10 * 60000).toISOString(),
    type: 'transcript',
    content: 'Let\'s implement the authentication system next.',
    speaker: 'Developer',
    duration: 8
  }
];

export const Default: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          transcript: sampleTranscript,
          loading: { ...streamInitialState.loading, transcript: false },
          error: { ...streamInitialState.error, transcript: null }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream/transcript`, () => {
          return HttpResponse.json(sampleTranscript);
        })
      ],
    },
  },
};

export const AdminView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: {
          ...streamInitialState,
          transcript: sampleTranscript,
          loading: { ...streamInitialState.loading, transcript: false },
          error: { ...streamInitialState.error, transcript: null }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream/transcript`, () => {
          return HttpResponse.json(sampleTranscript);
        })
      ],
    },
  },
};

export const Loading: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          loading: { ...streamInitialState.loading, transcript: true },
          error: { ...streamInitialState.error, transcript: null }
        },
      },
    },
  },
};

export const Error: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          loading: { ...streamInitialState.loading, transcript: false },
          error: { ...streamInitialState.error, transcript: 'Failed to load transcript data' }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream/transcript`, () => {
          return new HttpResponse(null, { status: 500, statusText: 'Server Error' });
        })
      ],
    },
  },
};

export const EmptyTranscript: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          transcript: [],
          loading: { ...streamInitialState.loading, transcript: false },
          error: { ...streamInitialState.error, transcript: null }
        },
      },
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/stream/transcript`, () => {
          return HttpResponse.json([]);
        })
      ],
    },
  },
};