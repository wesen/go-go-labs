import type { Meta, StoryObj } from '@storybook/react';
import React from 'react';
import StreamInfoDisplay from './StreamInfoDisplay';
import { store } from '../../.storybook/preview';
import {
  toggleEditMode,
  fetchStreamInfo,
  initialState as streamInitialState,
  setStreamInfo as setStreamInfoAction
} from '../store/slices/streamSlice';
import {
  setCredentials,
  logout as logoutAction,
  initialState as authInitialState,
  resetAuthState
} from '../store/slices/authSlice';

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

const resetStore = () => {
  store.dispatch({ type: 'STREAM_RESET_STATE', payload: { ...streamInitialState, info: baseStreamInfo } });
  store.dispatch(resetAuthState(authInitialState));
};

export const Default: Story = {
  play: async () => {
    resetStore();
    store.dispatch(setCredentials({ token: 'mock-token-admin', isAdmin: true }));
    store.dispatch(fetchStreamInfo.fulfilled(baseStreamInfo, "requestId", undefined));
  },
};

export const UserView: Story = {
  play: async () => {
    resetStore();
    store.dispatch(setCredentials({ token: 'mock-token-user', isAdmin: false }));
    store.dispatch(fetchStreamInfo.fulfilled(baseStreamInfo, "requestId", undefined));
  },
};

export const EditMode: Story = {
  play: async () => {
    resetStore();
    store.dispatch(setCredentials({ token: 'mock-token-admin', isAdmin: true }));
    store.dispatch(fetchStreamInfo.fulfilled(baseStreamInfo, "requestId", undefined));
    store.dispatch(toggleEditMode());
  },
};

export const Loading: Story = {
  play: async () => {
    resetStore();
    store.dispatch(setCredentials({ token: 'mock-token-admin', isAdmin: true }));
    store.dispatch(fetchStreamInfo.pending("requestId", undefined, undefined));
  },
};

export const Error: Story = {
  play: async () => {
    resetStore();
    store.dispatch(setCredentials({ token: 'mock-token-admin', isAdmin: true }));
    const errorPayload: Error = new Error('Failed to load stream information (Storybook mock)');
    store.dispatch(fetchStreamInfo.rejected(errorPayload, "requestId", undefined, undefined));
  },
};