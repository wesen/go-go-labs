import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../store';
import {
  fetchAllSteps,
  addNewUpcomingStep,
  completeActiveStep,
  reactivateStepFromSource
} from '../store/slices/streamSlice';
import { useGetAllStepsQuery } from '../api/stepsApi';
import {
  selectStepsData,
  selectStepsLoading,
  selectStepsError
} from '../store/selectors';

const TaskStepsPanel: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [newStep, setNewStep] = useState('');
  
  // Use memoized selectors to prevent unnecessary re-renders
  const { completedSteps, activeStep, upcomingSteps, stepIdMapping } = useSelector(selectStepsData);
  const loading = useSelector(selectStepsLoading);
  const error = useSelector(selectStepsError);
  
  // RTK Query approach
  const { isLoading: stepsQueryLoading, refetch: refetchSteps } = useGetAllStepsQuery(undefined, {
    // Skip initial fetch since we'll use the thunk for the first load
    skip: true
  });

  // Fetch steps on component mount
  useEffect(() => {
    dispatch(fetchAllSteps())
      .catch(error => {
        console.error('Failed to fetch steps:', error);
      });
  }, [dispatch]);

  const handleAddUpcomingStep = () => {
    if (newStep.trim()) {
      dispatch(addNewUpcomingStep(newStep.trim()))
        .unwrap()
        .then(() => setNewStep(''))
        .catch(error => {
          console.error('Failed to add upcoming step:', error);
        });
    }
  };

  const handleCompleteCurrentStep = () => {
    dispatch(completeActiveStep())
      .catch(error => {
        console.error('Failed to complete current step:', error);
      });
  };

  const handleSetActive = (step: string, source: 'upcoming' | 'completed') => {
    const stepId = stepIdMapping?.[step] || '';
    dispatch(reactivateStepFromSource({ step, source, stepId }))
      .catch(error => {
        console.error('Failed to reactivate step:', error);
      });
  };

  if (loading || stepsQueryLoading) {
    return <div className="p-4 bg-white rounded-lg shadow">Loading tasks...</div>;
  }

  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-800 rounded-lg shadow">
        Error loading tasks: {error}
        <button 
          onClick={() => dispatch(fetchAllSteps())} 
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Task Steps</h2>
      
      {/* Completed Steps */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2 text-gray-700">Completed</h3>
        {!completedSteps || completedSteps.length === 0 ? (
          <p className="text-gray-500 italic">No completed tasks yet</p>
        ) : (
          <ul className="space-y-2">
            {completedSteps.map((step) => (
              <li key={step} className="flex items-center">
                <span className="mr-2 text-green-500">✓</span>
                <span>{step}</span>
                <button
                  onClick={() => handleSetActive(step, 'completed')}
                  className="ml-auto text-xs bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                >
                  Reactivate
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
      
      {/* Active Step */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2 text-blue-600">Current Task</h3>
        {activeStep ? (
          <div className="bg-blue-50 p-3 rounded border border-blue-200 flex justify-between items-center">
            <span className="font-medium">{activeStep}</span>
            <button
              onClick={handleCompleteCurrentStep}
              className="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm"
            >
              Complete
            </button>
          </div>
        ) : (
          <p className="text-gray-500 italic">No active task</p>
        )}
      </div>
      
      {/* Upcoming Steps */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2 text-gray-700">Upcoming</h3>
        {!upcomingSteps || upcomingSteps.length === 0 ? (
          <p className="text-gray-500 italic">No upcoming tasks</p>
        ) : (
          <ul className="space-y-2">
            {upcomingSteps.map((step) => (
              <li key={step} className="flex items-center">
                <span className="mr-2 text-gray-400">○</span>
                <span>{step}</span>
                <button
                  onClick={() => handleSetActive(step, 'upcoming')}
                  className="ml-auto text-xs bg-blue-200 hover:bg-blue-300 px-2 py-1 rounded"
                >
                  Make Active
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
      
      {/* Add New Task */}
      <div>
        <h3 className="text-lg font-semibold mb-2 text-gray-700">Add New Task</h3>
        <div className="flex">
          <input
            type="text"
            value={newStep}
            onChange={(e) => setNewStep(e.target.value)}
            className="flex-grow p-2 border border-gray-300 rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Enter new task..."
          />
          <button
            onClick={handleAddUpcomingStep}
            disabled={!newStep.trim()}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-r disabled:bg-gray-300"
          >
            Add
          </button>
        </div>
      </div>
    </div>
  );
};

export default TaskStepsPanel;