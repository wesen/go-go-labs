import type { Meta, StoryObj } from '@storybook/react';
import React from 'react';
import TabsNavigation from './TabsNavigation';
import { initialState as streamInitialState } from '../store/slices/streamSlice';
import { initialState as authInitialState } from '../store/slices/authSlice';

const meta: Meta<typeof TabsNavigation> = {
  title: 'Components/TabsNavigation',
  component: TabsNavigation,
  parameters: {
    layout: 'fullscreen',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof TabsNavigation>;

export const MainTab: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          activeTab: 'main'
        }
      }
    }
  }
};

export const NotesTab: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          activeTab: 'notes'
        }
      }
    }
  }
};

export const SummaryTab: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          activeTab: 'summary'
        }
      }
    }
  }
};

export const RawTranscriptTab: Story = {
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          activeTab: 'raw'
        }
      }
    }
  }
};

export const TabInteraction: Story = {
  render: () => (
    <div className="container mx-auto p-4">
      <TabsNavigation />
      <div className="p-4 border border-gray-200 rounded">
        <p>Tab content would appear here based on selected tab.</p>
      </div>
    </div>
  ),
  parameters: {
    redux: {
      preloadedState: {
        auth: { ...authInitialState },
        stream: {
          ...streamInitialState,
          activeTab: 'main' // Default tab
        }
      }
    }
  }
};