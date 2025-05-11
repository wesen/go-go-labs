import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { AppDispatch } from './store';
import { fetchStreamInfo, fetchAllSteps, fetchTranscript, fetchGithubInfo } from './store/slices/streamSlice';
import StreamInfoDisplay from './components/StreamInfoDisplay';
import TaskStepsPanel from './components/TaskStepsPanel';
import TranscriptPanel from './components/TranscriptPanel';
import GithubInfoPanel from './components/GithubInfoPanel';
import AuthStatus from './components/AuthStatus';
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
        const promises = [
          dispatch(fetchStreamInfo()),
          dispatch(fetchAllSteps()),
          dispatch(fetchTranscript()),
          dispatch(fetchGithubInfo())
        ];
        
        // We use Promise.allSettled instead of Promise.all to continue even if some promises fail
        await Promise.allSettled(promises);
      } catch (error) {
        console.error('Error initializing app data:', error);
      } finally {
        // Always set data as initialized so app can render
        // The mock middleware will provide fallback data
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
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-gray-800">Stream Task Overview</h1>
          <ErrorBoundary>
            <AuthStatus />
          </ErrorBoundary>
        </div>
        
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