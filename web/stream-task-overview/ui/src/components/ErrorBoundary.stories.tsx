import type { Meta, StoryObj } from '@storybook/react';
import React from 'react';
import ErrorBoundary from './ErrorBoundary';

const meta: Meta<typeof ErrorBoundary> = {
  title: 'Components/ErrorBoundary',
  component: ErrorBoundary,
  parameters: {
    layout: 'centered',
  },
  tags: [],
};

export default meta;
type Story = StoryObj<typeof ErrorBoundary>;

// A component that throws an error for testing the ErrorBoundary
const BuggyComponent = () => {
  throw new Error('This is a test error');
  return <div>This will never render</div>;
};

// A component that renders normally without errors
const NormalComponent = () => {
  return <div className="p-4 bg-white rounded-lg shadow">This is a normal component that doesn't throw errors</div>;
};

// Custom fallback component
const CustomFallback = () => (
  <div className="p-4 bg-yellow-100 text-yellow-800 rounded-lg shadow">
    <h2 className="text-xl font-bold mb-2">Custom Error Message</h2>
    <p>Something went wrong with this component.</p>
  </div>
);

export const WithoutError: Story = {
  render: () => (
    <ErrorBoundary>
      <NormalComponent />
    </ErrorBoundary>
  ),
};

export const WithError: Story = {
  render: () => (
    <ErrorBoundary>
      <BuggyComponent />
    </ErrorBoundary>
  ),
  parameters: {
    // Disable console error in Storybook to avoid cluttering the console
    // This is just for this story since we're intentionally causing an error
    loki: { skip: true },
  },
};

export const WithCustomFallback: Story = {
  render: () => (
    <ErrorBoundary fallback={<CustomFallback />}>
      <BuggyComponent />
    </ErrorBoundary>
  ),
  parameters: {
    // Disable console error in Storybook to avoid cluttering the console
    loki: { skip: true },
  },
};

export const NestedErrorBoundaries: Story = {
  render: () => (
    <ErrorBoundary>
      <div className="p-4 bg-white rounded-lg shadow">
        <h3 className="font-bold mb-2">Outer component works fine</h3>
        <ErrorBoundary>
          <BuggyComponent />
        </ErrorBoundary>
        <div className="mt-4">This content should still render</div>
      </div>
    </ErrorBoundary>
  ),
  parameters: {
    loki: { skip: true },
  },
};