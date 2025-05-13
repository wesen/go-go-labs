import type { Meta, StoryObj } from '@storybook/react';
import GithubInfoPanel from './GithubInfoPanel';
import { http, HttpResponse } from 'msw';
import { baseUrl } from '../api/baseApi';
import { initialState as streamInitialState } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof GithubInfoPanel> = {
  title: 'Components/GithubInfoPanel',
  component: GithubInfoPanel,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof GithubInfoPanel>;

const connectedGithubInfo = {
  repoUrl: 'https://github.com/organization/repo-name',
  isConnected: true,
  token: 'mock-github-token',
  repoOwner: 'organization',
  repoName: 'repo-name',
  currentBranch: 'main',
  latestCommit: {
    message: 'Update components with new styling',
    author: 'Developer Name',
    hash: 'abc123def456',
    date: new Date().toISOString(),
    url: 'https://github.com/organization/repo-name/commit/abc123def456'
  }
};

const notConnectedGithubInfo = {
  repoUrl: 'https://github.com/organization/repo-name',
  isConnected: false,
  token: '',
  repoOwner: 'organization',
  repoName: 'repo-name',
  currentBranch: '',
  latestCommit: {
    message: '',
    author: '',
    hash: '',
    date: '',
    url: ''
  }
};

export const Connected: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          github: connectedGithubInfo,
          loading: { ...streamInitialState.loading, github: false },
          error: { ...streamInitialState.error, github: null }
        }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/github/info`, () => {
          return HttpResponse.json(connectedGithubInfo);
        })
      ],
    },
  },
};

export const NotConnected: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          github: notConnectedGithubInfo,
          loading: { ...streamInitialState.loading, github: false },
          error: { ...streamInitialState.error, github: null }
        }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/github/info`, () => {
          return HttpResponse.json(notConnectedGithubInfo);
        })
      ],
    },
  },
};

export const AdminNotConnected: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true, token: 'mock-token-admin' },
        stream: {
          ...streamInitialState,
          github: notConnectedGithubInfo,
          loading: { ...streamInitialState.loading, github: false },
          error: { ...streamInitialState.error, github: null }
        }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/github/info`, () => {
          return HttpResponse.json(notConnectedGithubInfo);
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
          github: notConnectedGithubInfo,
          loading: { ...streamInitialState.loading, github: true },
          error: { ...streamInitialState.error, github: null }
        }
      }
    }
  },
};

export const Error: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          github: {
            ...notConnectedGithubInfo,
            isConnected: true
          },
          loading: { ...streamInitialState.loading, github: false },
          error: { ...streamInitialState.error, github: 'Failed to load GitHub information' }
        }
      }
    },
    msw: {
      handlers: [
        http.get(`${baseUrl}/github/info`, () => {
          return new HttpResponse(null, { status: 500, statusText: 'Server Error' });
        })
      ],
    },
  },
};

// Add a story for missing GitHub state
export const MissingGithubState: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false, token: 'mock-token' },
        stream: {
          ...streamInitialState,
          // Deliberately omitting the github property
          github: undefined,
          loading: { ...streamInitialState.loading, github: false },
          error: { ...streamInitialState.error, github: null }
        }
      }
    }
  },
};