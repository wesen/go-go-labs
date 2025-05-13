import type { Meta, StoryObj } from '@storybook/react';
import RawTranscriptPanel from './RawTranscriptPanel';
import { initialState as streamInitialState, TranscriptEntry } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof RawTranscriptPanel> = {
  title: 'Components/RawTranscriptPanel',
  component: RawTranscriptPanel,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof RawTranscriptPanel>;

// Sample transcript entries
const sampleTranscriptEntries: TranscriptEntry[] = [
  {
    id: 't1',
    timestamp: new Date(Date.now() - 30 * 60000).toISOString(), // 30 minutes ago
    type: 'transcript',
    content: "Let's start by discussing our approach to the component library architecture.",
    speaker: "Host",
    duration: 12
  },
  {
    id: 't2',
    timestamp: new Date(Date.now() - 29 * 60000).toISOString(), // 29 minutes ago
    type: 'transcript',
    content: "I think we should follow atomic design principles. Starting with atoms, then building up to molecules and organisms.",
    speaker: "Guest",
    duration: 15
  },
  {
    id: 't3',
    timestamp: new Date(Date.now() - 28 * 60000).toISOString(), // 28 minutes ago
    type: 'transcript',
    content: "That makes sense. We'll need to define our design tokens first though, for colors, spacing, and typography.",
    speaker: "Host",
    duration: 10
  },
  {
    id: 't4',
    timestamp: new Date(Date.now() - 27 * 60000).toISOString(), // 27 minutes ago
    type: 'transcript',
    content: "Definitely. Let's also make sure we're building with accessibility in mind from the start.",
    speaker: "Guest",
    duration: 8
  },
  {
    id: 't5',
    timestamp: new Date(Date.now() - 25 * 60000).toISOString(), // 25 minutes ago
    type: 'transcript',
    content: "Absolutely! We should include keyboard navigation and proper ARIA attributes in all our components.",
    speaker: "Host",
    duration: 14
  },
];

// Mix of transcript and non-transcript entries
const mixedTranscriptEntries: TranscriptEntry[] = [
  ...sampleTranscriptEntries,
  {
    id: 'n1',
    timestamp: new Date(Date.now() - 20 * 60000).toISOString(), // 20 minutes ago
    type: 'note',
    content: 'Created initial design token structure'
  },
  {
    id: 'c1',
    timestamp: new Date(Date.now() - 15 * 60000).toISOString(), // 15 minutes ago
    type: 'commit',
    content: 'Initial component architecture setup',
    commitHash: 'a1b2c3d',
    commitUrl: 'https://github.com/org/repo/commit/a1b2c3d'
  },
];

export const WithTranscriptEntries: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false },
        stream: {
          ...streamInitialState,
          transcript: sampleTranscriptEntries
        }
      }
    }
  }
};

export const MixedEntryTypes: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false },
        stream: {
          ...streamInitialState,
          transcript: mixedTranscriptEntries
        }
      }
    }
  }
};

export const EmptyTranscript: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false },
        stream: {
          ...streamInitialState,
          transcript: []
        }
      }
    }
  }
};