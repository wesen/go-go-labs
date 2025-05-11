import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../store';
import { toggleEditMode, updateStreamInfo, fetchStreamInfo } from '../store/slices/streamSlice';
import { useGetStreamInfoQuery, useUpdateStreamInfoMutation } from '../api/streamApi';
import {
  selectStreamInfo,
  selectIsEditing,
  selectStreamInfoLoading,
  selectStreamInfoError
} from '../store/selectors';

interface EditableFieldProps {
  label: string;
  value: string;
  onChange: (value: string) => void;
  editMode: boolean;
}

const EditableField: React.FC<EditableFieldProps> = ({ label, value, onChange, editMode }) => {
  return (
    <div className="mb-4">
      <div className="text-sm text-gray-500 mb-1">{label}</div>
      {editMode ? (
        <input
          type="text"
          className="w-full p-2 border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={value}
          onChange={(e) => onChange(e.target.value)}
        />
      ) : (
        <div className="font-medium">{value}</div>
      )}
    </div>
  );
};

const StreamInfoDisplay: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  
  // Using RTK Query directly
  const { data: streamInfo, isLoading: isLoadingQuery, error: queryError } = useGetStreamInfoQuery();
  const [updateInfo, { isLoading: isUpdating }] = useUpdateStreamInfoMutation();
  
  // Use memoized selectors to prevent unnecessary re-renders
  const info = useSelector(selectStreamInfo);
  const isEditing = useSelector(selectIsEditing);
  const loading = useSelector(selectStreamInfoLoading);
  const error = useSelector(selectStreamInfoError);

  // Fetch stream info on component mount
  useEffect(() => {
    dispatch(fetchStreamInfo());
  }, [dispatch]);

  const handleToggleEditMode = () => {
    dispatch(toggleEditMode());
  };

  const handleSaveChanges = async () => {
    // Option 1: Using RTK Query directly
    try {
      await updateInfo(info).unwrap();
      dispatch(toggleEditMode());
    } catch (err) {
      console.error('Failed to update stream info:', err);
    }
    
    // Option 2: Using Redux thunk
    // dispatch(updateStreamInfo(info))
    //   .unwrap()
    //   .then(() => dispatch(toggleEditMode()))
    //   .catch((err) => console.error('Failed to update stream info:', err));
  };

  const handleChange = (field: string, value: string) => {
    dispatch(updateStreamInfo({ ...info, [field]: value }));
  };

  if (loading || isLoadingQuery) {
    return <div className="p-4 bg-white rounded-lg shadow">Loading stream information...</div>;
  }

  if (error || queryError) {
    return (
      <div className="p-4 bg-red-100 text-red-800 rounded-lg shadow">
        Error loading stream information: {error || 'Failed to load data'}
      </div>
    );
  }

  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-bold">Stream Information</h2>
        <button
          onClick={isEditing ? handleSaveChanges : handleToggleEditMode}
          className={`px-4 py-2 rounded ${isEditing ? 'bg-green-500 hover:bg-green-600' : 'bg-blue-500 hover:bg-blue-600'} text-white`}
          disabled={isUpdating}
        >
          {isUpdating ? 'Saving...' : isEditing ? 'Save Changes' : 'Edit'}
        </button>
      </div>

      <EditableField
        label="Title"
        value={info.title}
        onChange={(value) => handleChange('title', value)}
        editMode={isEditing}
      />

      <EditableField
        label="Description"
        value={info.description}
        onChange={(value) => handleChange('description', value)}
        editMode={isEditing}
      />

      <EditableField
        label="Programming Language"
        value={info.language}
        onChange={(value) => handleChange('language', value)}
        editMode={isEditing}
      />

      <EditableField
        label="GitHub Repository"
        value={info.githubRepo}
        onChange={(value) => handleChange('githubRepo', value)}
        editMode={isEditing}
      />

      <div className="mt-4 flex justify-between text-sm text-gray-500">
        <div>Started: {new Date(info.startTime).toLocaleString()}</div>
        <div>Viewers: {info.viewerCount}</div>
      </div>
    </div>
  );
};

export default StreamInfoDisplay;