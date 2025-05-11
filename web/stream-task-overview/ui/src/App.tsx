import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { AppDispatch } from './store';
import { fetchStreamInfo, fetchAllSteps, fetchTranscript, fetchGithubInfo } from './store/slices/streamSlice';
import StreamInfoDisplay from './components/StreamInfoDisplay';
import TaskStepsPanel from './components/TaskStepsPanel';
import TranscriptPanel from './components/TranscriptPanel';
import GithubInfoPanel from './components/GithubInfoPanel';
import ErrorBoundary from './components/ErrorBoundary';

const App: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [dataInitialized, setDataInitialized] = useState(false);

  // Initialize app data
  useEffect(() => {
    const initializeData = async () => {
      try {
        // Use try/catch to handle potential errors during initialization
        // Note: In a real app with a real backend, these would succeed
        // The current errors are because there's no actual backend running
        await dispatch(fetchStreamInfo()).unwrap();
        await dispatch(fetchAllSteps()).unwrap();
        await dispatch(fetchTranscript()).unwrap();
        await dispatch(fetchGithubInfo()).unwrap();
      } catch (error) {
        console.error('Error initializing app data:', error);
      } finally {
        setDataInitialized(true);
      }
    };

    initializeData();
  }, [dispatch]);

  // Show a loading state until data is initialized
  if (!dataInitialized) {
    return (
      <div className="min-h-screen bg-gray-100 p-6 flex items-center justify-center">
        <div className="text-xl font-medium text-gray-600">Loading application data...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <div className="max-w-7xl mx-auto">
        <h1 className="text-3xl font-bold mb-8 text-gray-800">Stream Task Overview</h1>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div className="md:col-span-2 space-y-6">
            <ErrorBoundary>
              <StreamInfoDisplay />
            </ErrorBoundary>
            
            <ErrorBoundary>
              <TaskStepsPanel />
            </ErrorBoundary>
            
            <ErrorBoundary>
              <TranscriptPanel />
            </ErrorBoundary>
          </div>
          
          <div>
            <ErrorBoundary>
              <GithubInfoPanel />
            </ErrorBoundary>
          </div>
        </div>
      </div>
    </div>
  );
};

export default App;