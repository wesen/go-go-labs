import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQuery, toCamelCase } from './baseApi';

// Define types for steps data
interface StepData {
  id: string;
  description: string;
  createdAt: string;
}

interface BackendStepInfo {
  completed: StepData[];
  active: StepData | null;
  upcoming: StepData[];
}

interface FrontendStepState {
  completedSteps: string[];
  activeStep: string;
  upcomingSteps: string[];
  _stepIdMapping: Record<string, string>; // Maps step description to id
}

// Transform backend step data to frontend format
export const toFrontendSteps = (backendSteps: BackendStepInfo): FrontendStepState => {
  // Create mapping of step description to id
  const idMapping: Record<string, string> = {};
  
  // Map completed steps
  const completedSteps = backendSteps.completed.map(step => {
    idMapping[step.description] = step.id;
    return step.description;
  });
  
  // Map active step
  let activeStep = '';
  if (backendSteps.active) {
    activeStep = backendSteps.active.description;
    idMapping[activeStep] = backendSteps.active.id;
  }
  
  // Map upcoming steps
  const upcomingSteps = backendSteps.upcoming.map(step => {
    idMapping[step.description] = step.id;
    return step.description;
  });
  
  return {
    completedSteps,
    activeStep,
    upcomingSteps,
    _stepIdMapping: idMapping
  };
};

// Helper to find step ID from name and source
export const getStepIdFromMapping = (state: FrontendStepState, step: string): string => {
  return state._stepIdMapping[step] || '';
};

export const stepsApi = createApi({
  reducerPath: 'stepsApi',
  baseQuery,
  tagTypes: ['Steps'],
  endpoints: (builder) => ({
    getAllSteps: builder.query<FrontendStepState, void>({
      query: () => 'stream/steps',
      transformResponse: (response) => toFrontendSteps(toCamelCase(response)),
      providesTags: ['Steps'],
    }),
    setActiveStep: builder.mutation<FrontendStepState, string>({
      query: (step) => ({
        url: 'stream/steps/active',
        method: 'PUT',
        body: { step }
      }),
      transformResponse: (response) => toFrontendSteps(toCamelCase(response)),
      invalidatesTags: ['Steps'],
    }),
    addUpcomingStep: builder.mutation<FrontendStepState, string>({
      query: (step) => ({
        url: 'stream/steps/upcoming',
        method: 'POST',
        body: { step }
      }),
      transformResponse: (response) => toFrontendSteps(toCamelCase(response)),
      invalidatesTags: ['Steps'],
    }),
    completeCurrentStep: builder.mutation<FrontendStepState, void>({
      query: () => ({
        url: 'stream/steps/complete',
        method: 'POST'
      }),
      transformResponse: (response) => toFrontendSteps(toCamelCase(response)),
      invalidatesTags: ['Steps'],
    }),
    reactivateStep: builder.mutation<FrontendStepState, {step: string, source: 'upcoming' | 'completed', stepId: string}>({
      query: ({stepId, source}) => ({
        url: 'stream/steps/reactivate',
        method: 'PUT',
        body: { stepId, source }
      }),
      transformResponse: (response) => toFrontendSteps(toCamelCase(response)),
      invalidatesTags: ['Steps'],
    }),
  }),
});

export const {
  useGetAllStepsQuery,
  useSetActiveStepMutation,
  useAddUpcomingStepMutation,
  useCompleteCurrentStepMutation,
  useReactivateStepMutation
} = stepsApi;