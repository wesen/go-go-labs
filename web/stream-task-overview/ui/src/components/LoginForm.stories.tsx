import type { Meta, StoryObj } from '@storybook/react';
import { Provider } from 'react-redux';
import LoginForm from './LoginForm';
import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../store/slices/authSlice';
import { authApi } from '../api/authApi';

const meta: Meta<typeof LoginForm> = {
  title: 'Components/LoginForm',
  component: LoginForm,
  decorators: [
    (Story) => (
      <Provider store={configureStore({
        reducer: {
          auth: authReducer,
          [authApi.reducerPath]: authApi.reducer,
        },
        middleware: (getDefaultMiddleware) =>
          getDefaultMiddleware().concat(authApi.middleware),
      })}>
        <div className="max-w-md mx-auto p-4">
          <Story />
        </div>
      </Provider>
    ),
  ],
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
};

export default meta;
type Story = StoryObj<typeof LoginForm>;

export const Default: Story = {};

export const WithCallback: Story = {
  args: {
    onSuccess: () => alert('Login successful!')
  }
};