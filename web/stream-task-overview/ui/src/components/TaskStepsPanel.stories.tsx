import type { Meta, StoryObj } from '@storybook/react';
import TaskStepsPanel from './TaskStepsPanel';
import { initialState as streamInitialState } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof TaskStepsPanel> = {
  title: 'Components/TaskStepsPanel',
  component: TaskStepsPanel,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof TaskStepsPanel>;

export const Default: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true },
        stream: {
          ...streamInitialState,
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
            ...streamInitialState.loading,
            steps: false
          },
          error: {
            ...streamInitialState.error,
            steps: null
          }
        }
      }
    }
  }
};

export const Loading: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true },
        stream: {
          ...streamInitialState,
          completedSteps: [],
          activeStep: '',
          upcomingSteps: [],
          stepIdMapping: {},
          loading: {
            ...streamInitialState.loading,
            steps: true
          },
          error: {
            ...streamInitialState.error,
            steps: null
          }
        }
      }
    }
  }
};

export const Error: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true },
        stream: {
          ...streamInitialState,
          completedSteps: [],
          activeStep: '',
          upcomingSteps: [],
          stepIdMapping: {},
          loading: {
            ...streamInitialState.loading,
            steps: false
          },
          error: {
            ...streamInitialState.error,
            steps: 'Failed to load steps'
          }
        }
      }
    }
  }
};

// Admin view is the default
export const AdminView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: true },
        stream: {
          ...streamInitialState,
          completedSteps: ['Research competitors', 'Create wireframes'],
          activeStep: 'Implement UI components',
          upcomingSteps: ['Write unit tests', 'Deploy to staging'],
          stepIdMapping: {
            'Research competitors': '1',
            'Create wireframes': '2',
            'Implement UI components': '3',
            'Write unit tests': '4',
            'Deploy to staging': '5'
          }
        }
      }
    }
  }
};

// Regular user view
export const UserView: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState, isAuthenticated: true, isAdmin: false },
        stream: {
          ...streamInitialState,
          completedSteps: ['Research competitors', 'Create wireframes'],
          activeStep: 'Implement UI components',
          upcomingSteps: ['Write unit tests', 'Deploy to staging'],
          stepIdMapping: {
            'Research competitors': '1',
            'Create wireframes': '2',
            'Implement UI components': '3',
            'Write unit tests': '4',
            'Deploy to staging': '5'
          }
        }
      }
    }
  }
};