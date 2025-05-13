import type { Meta, StoryObj } from '@storybook/react';
import { EditableField } from './StreamInfoDisplay'; // This component may need to be exported separately

const meta: Meta<typeof EditableField> = {
  title: 'Components/EditableField',
  component: EditableField,
  tags: [],
  argTypes: {
    label: { control: 'text' },
    value: { control: 'text' },
    onChange: { action: 'changed' },
    editMode: { control: 'boolean' }
  }
};

export default meta;
type Story = StoryObj<typeof EditableField>;

export const ReadOnly: Story = {
  args: {
    label: 'Field Label',
    value: 'Field Value',
    editMode: false
  }
};

export const Editing: Story = {
  args: {
    label: 'Field Label',
    value: 'Field Value',
    editMode: true
  }
};

export const LongText: Story = {
  args: {
    label: 'Description',
    value: 'This is a very long text that should demonstrate how the component handles longer content in both read and edit modes.',
    editMode: false
  }
};

export const EditingLongText: Story = {
  args: {
    label: 'Description',
    value: 'This is a very long text that should demonstrate how the component handles longer content in both read and edit modes.',
    editMode: true
  }
};