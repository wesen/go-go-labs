import type { Meta, StoryObj } from '@storybook/react';
import { Provider } from 'react-redux';
import TaskStepsPanel from './TaskStepsPanel';
import { createMockStore } from '../mocks/store';

// We'll use MSW to mock the auth API response instead of mocking the hook directly

// Create mock stores with proper data structure
const mockStore = createMockStore({
  stream: {
    completedSteps: ['Research competitors', 'Create wireframes'],
    activeStep: 'Implement UI components',
    upcomingSteps: ['Write unit tests', 'Deploy to staging'],
    stepIdMapping: {
      'Research competitors': '1',
      'Create wireframes': '2',
      'Implement UI components': '3',
      'Write unit tests': '4',
      'Deploy to staging': '5'
    },
    loading: {
      streamInfo: false,
      steps: false,
      transcript: false,
      github: false
    },
    error: {
      streamInfo: null,
      steps: null,
      transcript: null,
      github: null
    }
  }
});

const meta: Meta<typeof TaskStepsPanel> = {
  title: 'Components/TaskStepsPanel',
  component: TaskStepsPanel,
  decorators: [
    (Story) => (
      <Provider store={mockStore}>
        <Story />
      </Provider>
    ),
  ],
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof TaskStepsPanel>;

export const Default: Story = {};

export const Loading: Story = {
  decorators: [
    (Story) => (
      <Provider store={createMockStore({
        stream: {
          completedSteps: [],
          activeStep: '',
          upcomingSteps: [],
          stepIdMapping: {},
          loading: {
            streamInfo: false,
            steps: true,
            transcript: false,
            github: false
          },
          error: {
            streamInfo: null,
            steps: null,
            transcript: null,
            github: null
          }
        }
      })}>
        <Story />
      </Provider>
    ),
  ],
};

export const Error: Story = {
  decorators: [
    (Story) => (
      <Provider store={createMockStore({
        stream: {
          completedSteps: [],
          activeStep: '',
          upcomingSteps: [],
          stepIdMapping: {},
          loading: {
            streamInfo: false,
            steps: false,
            transcript: false,
            github: false
          },
          error: {
            streamInfo: null,
            steps: 'Failed to load steps',
            transcript: null,
            github: null
          }
        }
      })}>
        <Story />
      </Provider>
    ),
  ],
};

// No need for special configuration - default is admin view
export const AdminView: Story = {};

// Create a custom version for non-admin users by overriding store state
export const UserView: Story = {
  parameters: {
    mockData: {
      isAdmin: false
    }
  }
};