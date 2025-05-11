import { isRejectedWithValue } from '@reduxjs/toolkit'
import type { Middleware } from '@reduxjs/toolkit'
import { initialState } from '../store/slices/streamSlice';

/**
 * Mock middleware to handle API requests when the backend is not available
 * This is useful for development when you don't have a backend running
 * or when there are CORS issues that prevent the frontend from connecting to the backend
 */
export const mockApiMiddleware: Middleware = api => next => action => {
  // Check if it's a rejected API call with CORS error
  if (isRejectedWithValue(action)) {
    const error = action.payload?.error || action.error;
    const errorMessage = error?.message || '';
    
    // Check if this is a CORS error
    const isCorsError = errorMessage.includes('CORS') || 
                       errorMessage.includes('Failed to fetch') || 
                       errorMessage.includes('Network Error');
    
    if (isCorsError) {
      console.warn(
        'API request failed with possible CORS issue. Using mock data. In a production environment, make sure your API supports CORS.',
        action
      );
      
      // Create a mock response based on the endpoint and query
      const endpoint = action.meta?.arg?.endpointName || '';
      const mockResponse = createMockResponse(endpoint, action.meta?.arg?.originalArgs);
      
      // If we have a mock response, use it
      if (mockResponse) {
        // Dispatch a success action with our mock data
        return next({
          ...action,
          payload: mockResponse,
          meta: {
            ...action.meta,
            requestStatus: 'fulfilled'
          },
          type: action.type.replace('/rejected', '/fulfilled')
        });
      }
    }
  }
  
  return next(action);
};

/**
 * Create a mock response based on the endpoint and arguments
 */
function createMockResponse(endpoint: string, args: any) {
  // Create mock responses for different endpoints
  switch (endpoint) {
    case 'getStreamInfo':
      return initialState.info;
      
    case 'getAllSteps':
      return {
        completedSteps: initialState.completedSteps,
        activeStep: initialState.activeStep,
        upcomingSteps: initialState.upcomingSteps,
        _stepIdMapping: initialState.stepIdMapping || {}
      };
      
    case 'getTranscript':
      return initialState.transcript;
      
    case 'getGithubInfo':
      return initialState.github;
      
    case 'addNote':
      return {
        id: Math.random().toString(36).substr(2, 9),
        timestamp: new Date().toISOString(),
        type: 'note',
        content: args || 'Mock note'
      };
      
    case 'addParagraph':
      return {
        id: Math.random().toString(36).substr(2, 9),
        timestamp: new Date().toISOString(),
        type: 'paragraph',
        content: args?.content || 'Mock paragraph',
        title: args?.title || 'Mock Title',
        timeRange: args?.timeRange || { 
          start: new Date(Date.now() - 10 * 60000).toISOString(),
          end: new Date().toISOString()
        }
      };
      
    case 'updateStreamInfo':
      return { ...initialState.info, ...args };
      
    case 'completeCurrentStep':
    case 'setActiveStep':
    case 'addUpcomingStep':
    case 'reactivateStep':
      return {
        completedSteps: initialState.completedSteps,
        activeStep: initialState.activeStep,
        upcomingSteps: initialState.upcomingSteps,
        _stepIdMapping: initialState.stepIdMapping || {}
      };
      
    default:
      return null;
  }
}