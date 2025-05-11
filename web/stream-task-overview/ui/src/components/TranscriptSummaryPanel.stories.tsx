import type { Meta, StoryObj } from '@storybook/react';
import TranscriptSummaryPanel from './TranscriptSummaryPanel';
import { initialState as streamInitialState, TranscriptEntry } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof TranscriptSummaryPanel> = {
  title: 'Components/TranscriptSummaryPanel',
  component: TranscriptSummaryPanel,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof TranscriptSummaryPanel>;

// Generate comprehensive sample data for a full stream session
const fullStreamData: TranscriptEntry[] = [
  // Task starts
  {
    id: 'ts1',
    timestamp: new Date(Date.now() - 90 * 60000).toISOString(), // 90 minutes ago
    type: 'task_started',
    content: 'Started working on project setup',
    taskName: 'Project Setup'
  },
  {
    id: 'ts2',
    timestamp: new Date(Date.now() - 60 * 60000).toISOString(), // 60 minutes ago
    type: 'task_started',
    content: 'Started working on component architecture',
    taskName: 'Component Architecture'
  },
  {
    id: 'ts3',
    timestamp: new Date(Date.now() - 30 * 60000).toISOString(), // 30 minutes ago
    type: 'task_started',
    content: 'Started working on styling system',
    taskName: 'Styling System'
  },
  
  // Task completions
  {
    id: 'tc1',
    timestamp: new Date(Date.now() - 65 * 60000).toISOString(), // 65 minutes ago
    type: 'task_completed',
    content: 'Completed project setup',
    taskName: 'Project Setup'
  },
  {
    id: 'tc2',
    timestamp: new Date(Date.now() - 35 * 60000).toISOString(), // 35 minutes ago
    type: 'task_completed',
    content: 'Completed component architecture',
    taskName: 'Component Architecture'
  },
  
  // Commits
  {
    id: 'c1',
    timestamp: new Date(Date.now() - 85 * 60000).toISOString(), // 85 minutes ago
    type: 'commit',
    content: 'Initial project setup with TypeScript',
    commitHash: 'a1b2c3d4e5f6',
    commitUrl: 'https://github.com/organization/repo/commit/a1b2c3d4e5f6'
  },
  {
    id: 'c2',
    timestamp: new Date(Date.now() - 75 * 60000).toISOString(), // 75 minutes ago
    type: 'commit',
    content: 'Add build configuration',
    commitHash: 'b2c3d4e5f6g7',
    commitUrl: 'https://github.com/organization/repo/commit/b2c3d4e5f6g7'
  },
  {
    id: 'c3',
    timestamp: new Date(Date.now() - 50 * 60000).toISOString(), // 50 minutes ago
    type: 'commit',
    content: 'Create component folder structure',
    commitHash: 'c3d4e5f6g7h8',
    commitUrl: 'https://github.com/organization/repo/commit/c3d4e5f6g7h8'
  },
  {
    id: 'c4',
    timestamp: new Date(Date.now() - 40 * 60000).toISOString(), // 40 minutes ago
    type: 'commit',
    content: 'Implement base component interfaces',
    commitHash: 'd4e5f6g7h8i9',
    commitUrl: 'https://github.com/organization/repo/commit/d4e5f6g7h8i9'
  },
  {
    id: 'c5',
    timestamp: new Date(Date.now() - 20 * 60000).toISOString(), // 20 minutes ago
    type: 'commit',
    content: 'Add TailwindCSS configuration',
    commitHash: 'e5f6g7h8i9j0',
    commitUrl: 'https://github.com/organization/repo/commit/e5f6g7h8i9j0'
  },
  
  // Notes
  {
    id: 'n1',
    timestamp: new Date(Date.now() - 80 * 60000).toISOString(), // 80 minutes ago
    type: 'note',
    content: 'Decided to use TypeScript for better type safety'
  },
  {
    id: 'n2',
    timestamp: new Date(Date.now() - 70 * 60000).toISOString(), // 70 minutes ago
    type: 'note',
    content: 'Using Vite as build tool for faster development experience'
  },
  {
    id: 'n3',
    timestamp: new Date(Date.now() - 55 * 60000).toISOString(), // 55 minutes ago
    type: 'note',
    content: 'Component architecture will follow atomic design principles'
  },
  {
    id: 'n4',
    timestamp: new Date(Date.now() - 45 * 60000).toISOString(), // 45 minutes ago
    type: 'note',
    content: 'Need to implement proper accessibility standards in all components'
  },
  {
    id: 'n5',
    timestamp: new Date(Date.now() - 25 * 60000).toISOString(), // 25 minutes ago
    type: 'note',
    content: 'Planning to use CSS modules with TailwindCSS utility classes'
  },
  {
    id: 'n6',
    timestamp: new Date(Date.now() - 15 * 60000).toISOString(), // 15 minutes ago
    type: 'note',
    content: 'Dark mode support will be implemented using CSS variables'
  },
];

// GitHub info for the summary
const githubInfo = {
  repoUrl: 'https://github.com/organization/repo',
  isConnected: true,
  token: 'mock-github-token',
  repoOwner: 'organization',
  repoName: 'repo',
  currentBranch: 'main',
  latestCommit: {
    message: 'Add TailwindCSS configuration',
    author: 'Developer Name',
    hash: 'e5f6g7h8i9j0',
    date: new Date(Date.now() - 20 * 60000).toISOString(),
    url: 'https://github.com/organization/repo/commit/e5f6g7h8i9j0'
  }
};

export const CompleteSummary: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          transcript: fullStreamData,
          github: githubInfo
        }
      }
    }
  }
};

export const NoCommits: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          transcript: fullStreamData.filter(entry => entry.type !== 'commit'),
          github: {
            ...githubInfo,
            isConnected: false,
            latestCommit: {
              message: '',
              author: '',
              hash: '',
              date: '',
              url: ''
            }
          }
        }
      }
    }
  }
};

export const OnlyTasks: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          transcript: fullStreamData.filter(
            entry => entry.type === 'task_started' || entry.type === 'task_completed'
          ),
          github: {
            ...githubInfo,
            isConnected: false
          }
        }
      }
    }
  }
};

export const EmptyTranscript: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          transcript: [],
          github: {
            ...githubInfo,
            isConnected: false
          }
        }
      }
    }
  }
};