import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface GithubInfo {
  repoUrl: string;
  isConnected: boolean;
  token: string;
  repoOwner: string;
  repoName: string;
  currentBranch: string;
  latestCommit: {
    message: string;
    author: string;
    hash: string;
    date: string;
    url: string;
  };
  error?: string;
}

interface StreamInfo {
  title: string;
  description: string;
  startTime: string;
  language: string;
  githubRepo: string;
  currentTask: string;
  viewerCount: number;
}

interface StreamState {
  info: StreamInfo;
  completedSteps: string[];
  activeStep: string;
  upcomingSteps: string[];
  isEditing: boolean;
  isLoggedIn: boolean;
  github: GithubInfo;
}

const initialState: StreamState = {
  info: {
    title: "Building a React Component Library",
    description: "Creating reusable UI components with TailwindCSS",
    startTime: new Date().toISOString(),
    language: "JavaScript/React",
    githubRepo: "https://github.com/yourusername/component-library",
    currentTask: "",
    viewerCount: 42,
  },
  completedSteps: [
    "Project setup and initialization",
    "Design system planning"
  ],
  activeStep: "Setting up component architecture",
  upcomingSteps: [
    "Implement Button component",
    "Create Card component",
    "Build Form elements",
    "Add dark mode toggle"
  ],
  isEditing: false,
  isLoggedIn: true,
  github: {
    repoUrl: "https://github.com/yourusername/component-library",
    isConnected: false,
    token: "",
    repoOwner: "yourusername",
    repoName: "component-library",
    currentBranch: "main",
    latestCommit: {
      message: "",
      author: "",
      hash: "",
      date: "",
      url: ""
    }
  }
};

export const streamSlice = createSlice({
  name: 'stream',
  initialState,
  reducers: {
    setStreamInfo: (state, action: PayloadAction<StreamInfo>) => {
      state.info = action.payload;
      
      // If GitHub repo URL changed, parse owner and repo name
      if (action.payload.githubRepo !== state.info.githubRepo) {
        const repoUrlMatch = action.payload.githubRepo.match(/github\.com[\/:]([\w-\.]+)\/([\w-\.]+)(\.git)?$/);
        if (repoUrlMatch) {
          state.github.repoUrl = action.payload.githubRepo;
          state.github.repoOwner = repoUrlMatch[1];
          state.github.repoName = repoUrlMatch[2];
        }
      }
    },
    toggleEditMode: (state) => {
      state.isEditing = !state.isEditing;
    },
    toggleLoggedIn: (state) => {
      state.isLoggedIn = !state.isLoggedIn;
      // If logging out, cancel any editing mode
      if (!state.isLoggedIn) {
        state.isEditing = false;
      }
    },
    resetTimer: (state) => {
      state.info.startTime = new Date().toISOString();
    },
    addUpcomingStep: (state, action: PayloadAction<string>) => {
      state.upcomingSteps.push(action.payload);
    },
    setNewActiveTopic: (state, action: PayloadAction<string>) => {
      if (state.activeStep) {
        state.completedSteps.push(state.activeStep);
      }
      state.activeStep = action.payload;
    },
    completeCurrentStep: (state) => {
      if (state.activeStep) {
        state.completedSteps.push(state.activeStep);
        if (state.upcomingSteps.length > 0) {
          state.activeStep = state.upcomingSteps[0];
          state.upcomingSteps.splice(0, 1);
        } else {
          state.activeStep = "";
        }
      }
    },
    makeStepActive: (state, action: PayloadAction<{step: string, source: 'upcoming' | 'completed'}>) => {
      const { step, source } = action.payload;
      
      if (state.activeStep) {
        state.completedSteps.push(state.activeStep);
      }
      
      state.activeStep = step;
      
      if (source === 'upcoming') {
        state.upcomingSteps = state.upcomingSteps.filter(s => s !== step);
      } else if (source === 'completed') {
        state.completedSteps = state.completedSteps.filter(s => s !== step);
      }
    },
    // GitHub related actions
    setGithubToken: (state, action: PayloadAction<string>) => {
      state.github.token = action.payload;
    },
    setGithubConnectionStatus: (state, action: PayloadAction<boolean>) => {
      state.github.isConnected = action.payload;
    },
    setGithubError: (state, action: PayloadAction<string>) => {
      state.github.error = action.payload;
    },
    updateGithubBranchInfo: (state, action: PayloadAction<{branch: string}>) => {
      state.github.currentBranch = action.payload.branch;
    },
    updateGithubCommitInfo: (state, action: PayloadAction<{
      message: string;
      author: string;
      hash: string;
      date: string;
      url: string;
    }>) => {
      state.github.latestCommit = action.payload;
    }
  }
});

export const { 
  setStreamInfo, 
  toggleEditMode, 
  toggleLoggedIn,
  resetTimer, 
  addUpcomingStep, 
  setNewActiveTopic, 
  completeCurrentStep, 
  makeStepActive,
  // GitHub actions
  setGithubToken,
  setGithubConnectionStatus,
  setGithubError,
  updateGithubBranchInfo,
  updateGithubCommitInfo
} = streamSlice.actions;

export default streamSlice.reducer;