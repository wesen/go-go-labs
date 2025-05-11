import { createSlice, PayloadAction } from '@reduxjs/toolkit';

interface TranscriptEntry {
  id: string;
  timestamp: string;
  type: 'task_started' | 'task_completed' | 'commit' | 'note';
  content: string;
  taskName?: string;
  commitHash?: string;
  commitUrl?: string;
}

type TabType = 'main' | 'transcript' | 'summary';

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
  transcript: TranscriptEntry[];
  activeTab: TabType;
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
  },
  // Mock transcript data
  transcript: [
    {
      id: '1',
      timestamp: new Date(Date.now() - 45 * 60000).toISOString(), // 45 minutes ago
      type: 'task_started',
      content: 'Started working on project setup and initialization',
      taskName: 'Project setup and initialization'
    },
    {
      id: '2',
      timestamp: new Date(Date.now() - 40 * 60000).toISOString(), // 40 minutes ago
      type: 'note',
      content: 'Created project structure using create-react-app with TypeScript template'
    },
    {
      id: '3',
      timestamp: new Date(Date.now() - 35 * 60000).toISOString(), // 35 minutes ago
      type: 'commit',
      content: 'Initial project setup with TypeScript configuration',
      commitHash: 'a1b2c3d',
      commitUrl: 'https://github.com/yourusername/component-library/commit/a1b2c3d'
    },
    {
      id: '4',
      timestamp: new Date(Date.now() - 30 * 60000).toISOString(), // 30 minutes ago
      type: 'task_completed',
      content: 'Completed project setup and initialization',
      taskName: 'Project setup and initialization'
    },
    {
      id: '5',
      timestamp: new Date(Date.now() - 25 * 60000).toISOString(), // 25 minutes ago
      type: 'task_started',
      content: 'Started working on design system planning',
      taskName: 'Design system planning'
    },
    {
      id: '6',
      timestamp: new Date(Date.now() - 20 * 60000).toISOString(), // 20 minutes ago
      type: 'note',
      content: 'Researching color palette and typography options for the design system'
    },
    {
      id: '7',
      timestamp: new Date(Date.now() - 15 * 60000).toISOString(), // 15 minutes ago
      type: 'commit',
      content: 'Add design system tokens and Tailwind configuration',
      commitHash: 'e5f6g7h',
      commitUrl: 'https://github.com/yourusername/component-library/commit/e5f6g7h'
    },
    {
      id: '8',
      timestamp: new Date(Date.now() - 10 * 60000).toISOString(), // 10 minutes ago
      type: 'task_completed',
      content: 'Completed design system planning',
      taskName: 'Design system planning'
    },
    {
      id: '9',
      timestamp: new Date(Date.now() - 5 * 60000).toISOString(), // 5 minutes ago
      type: 'task_started',
      content: 'Started working on component architecture',
      taskName: 'Setting up component architecture'
    },
    {
      id: '10',
      timestamp: new Date().toISOString(), // now
      type: 'note',
      content: 'Creating folder structure for components and defining TypeScript interfaces'
    }
  ],
  activeTab: 'main'
};

// Helper function to generate a transcript entry
const generateTranscriptEntry = (type: TranscriptEntry['type'], content: string, details?: Partial<TranscriptEntry>): TranscriptEntry => ({
  id: Math.random().toString(36).substr(2, 9),
  timestamp: new Date().toISOString(),
  type,
  content,
  ...details
});

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
        // Add task completion to transcript
        state.transcript.push(generateTranscriptEntry('task_completed', 
          `Completed ${state.activeStep}`, { taskName: state.activeStep }));
      }
      state.activeStep = action.payload;
      // Add task start to transcript
      state.transcript.push(generateTranscriptEntry('task_started', 
        `Started working on ${action.payload}`, { taskName: action.payload }));
    },
    completeCurrentStep: (state) => {
      if (state.activeStep) {
        const completedTask = state.activeStep;
        state.completedSteps.push(completedTask);
        
        // Add task completion to transcript
        state.transcript.push(generateTranscriptEntry('task_completed', 
          `Completed ${completedTask}`, { taskName: completedTask }));
        
        if (state.upcomingSteps.length > 0) {
          state.activeStep = state.upcomingSteps[0];
          state.upcomingSteps.splice(0, 1);
          
          // Add new task start to transcript
          state.transcript.push(generateTranscriptEntry('task_started', 
            `Started working on ${state.activeStep}`, { taskName: state.activeStep }));
        } else {
          state.activeStep = "";
        }
      }
    },
    makeStepActive: (state, action: PayloadAction<{step: string, source: 'upcoming' | 'completed'}>) => {
      const { step, source } = action.payload;
      
      if (state.activeStep) {
        const previousTask = state.activeStep;
        state.completedSteps.push(previousTask);
        
        // Add task completion to transcript
        state.transcript.push(generateTranscriptEntry('task_completed', 
          `Completed ${previousTask}`, { taskName: previousTask }));
      }
      
      state.activeStep = step;
      
      // Add new task start to transcript
      state.transcript.push(generateTranscriptEntry('task_started', 
        `Started working on ${step}`, { taskName: step }));
      
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
      
      // Add commit to transcript
      state.transcript.push(generateTranscriptEntry('commit', 
        action.payload.message, {
          commitHash: action.payload.hash,
          commitUrl: action.payload.url
        }));
    },
    
    // Transcript related actions
    changeTab: (state, action: PayloadAction<TabType>) => {
      state.activeTab = action.payload;
    },
    addTranscriptNote: (state, action: PayloadAction<string>) => {
      state.transcript.push(generateTranscriptEntry('note', action.payload));
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
  updateGithubCommitInfo,
  // Transcript actions
  changeTab,
  addTranscriptNote
} = streamSlice.actions;

export default streamSlice.reducer;