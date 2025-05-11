import { createSelector } from '@reduxjs/toolkit';
import { RootState } from './index';

// Base selectors
const selectStreamState = (state: RootState) => state.stream;
const selectAuthState = (state: RootState) => state.auth;

// Auth selectors
export const selectIsAdmin = createSelector(
  selectAuthState,
  (authState) => authState.isAdmin
);

export const selectIsAuthenticated = createSelector(
  selectAuthState,
  (authState) => authState.isAuthenticated
);

// Memoized selectors for StreamInfoDisplay
export const selectStreamInfo = createSelector(
  selectStreamState,
  (streamState) => streamState.info
);

export const selectIsEditing = createSelector(
  selectStreamState,
  (streamState) => streamState.isEditing
);

export const selectStreamInfoLoading = createSelector(
  selectStreamState,
  (streamState) => streamState.loading.streamInfo
);

export const selectStreamInfoError = createSelector(
  selectStreamState,
  (streamState) => streamState.error.streamInfo
);

// Memoized selectors for TaskStepsPanel
export const selectStepsData = createSelector(
  selectStreamState,
  (streamState) => ({
    completedSteps: streamState.completedSteps,
    activeStep: streamState.activeStep,
    upcomingSteps: streamState.upcomingSteps,
    stepIdMapping: streamState.stepIdMapping
  })
);

export const selectStepsLoading = createSelector(
  selectStreamState,
  (streamState) => streamState.loading.steps
);

export const selectStepsError = createSelector(
  selectStreamState,
  (streamState) => streamState.error.steps
);

// Memoized selectors for TranscriptPanel
export const selectTranscriptData = createSelector(
  selectStreamState,
  (streamState) => streamState.transcript || []
);

export const selectTranscriptLoading = createSelector(
  selectStreamState,
  (streamState) => streamState.loading.transcript
);

export const selectTranscriptError = createSelector(
  selectStreamState,
  (streamState) => streamState.error.transcript
);

// Memoized selectors for GithubInfoPanel
export const selectGithubData = createSelector(
  selectStreamState,
  (streamState) => streamState.github
);

export const selectGithubLoading = createSelector(
  selectStreamState,
  (streamState) => streamState.loading.github
);

export const selectGithubError = createSelector(
  selectStreamState,
  (streamState) => streamState.error.github
);