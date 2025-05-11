import { createSlice, createAsyncThunk, PayloadAction } from '@reduxjs/toolkit';
import { streamApi } from '../../api/streamApi';
import { stepsApi } from '../../api/stepsApi';
import { transcriptApi } from '../../api/transcriptApi';
import { githubApi } from '../../api/githubApi';

export interface TranscriptEntry {
  id: string;
  timestamp: string;
  type: 'task_started' | 'task_completed' | 'commit' | 'note' | 'paragraph' | 'transcript';
  content: string;
  taskName?: string;
  commitHash?: string;
  commitUrl?: string;
  // For LLM-generated paragraphs
  timeRange?: { start: string; end: string };
  title?: string;
  // For transcript entries
  speaker?: string;
  duration?: number; // in seconds
}

type TabType = 'main' | 'notes' | 'summary' | 'raw';

export interface GithubInfo {
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

export interface StreamInfo {
  title: string;
  description: string;
  startTime: string;
  language: string;
  githubRepo: string;
  currentTask: string;
  viewerCount: number;
}

export interface StreamState {
  info: StreamInfo;
  completedSteps: string[];
  activeStep: string;
  upcomingSteps: string[];
  isEditing: boolean;
  isLoggedIn: boolean;
  github: GithubInfo;
  transcript: TranscriptEntry[];
  activeTab: TabType;
  // API state
  loading: {
    streamInfo: boolean;
    steps: boolean;
    transcript: boolean;
    github: boolean;
  };
  error: {
    streamInfo: string | null;
    steps: string | null;
    transcript: string | null;
    github: string | null;
  };
  // For step operations that use IDs
  stepIdMapping: Record<string, string>;
}

export const initialState: StreamState = {
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
  // Initialize API state
  loading: {
    streamInfo: false,
    steps: false,
    transcript: false,
    github: false
  },
  error: {
    streamInfo: null,
    steps: null,
    transcript: null,
    github: null
  },
  stepIdMapping: {},
  // Mock transcript/notes data
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
    },
    {
      id: '11',
      timestamp: new Date(Date.now() - 28 * 60000).toISOString(), // 28 minutes ago
      type: 'paragraph',
      title: 'Project Initialization Discussion',
      content: "During the initial setup phase, we discussed several approaches to structuring the React component library. We decided to use TypeScript for type safety and better developer experience, along with a modular architecture that separates component logic from styling.\n\nWe explored different build tools and settled on Vite for its speed and modern features. The decision to use CSS modules rather than styled-components was made after considering bundle size and performance implications.\n\nKey points discussed:\n- TypeScript configuration with strict mode enabled\n- Component directory structure following atomic design principles\n- Test strategy using Jest and React Testing Library\n- Documentation approach with Storybook",
      timeRange: {
        start: new Date(Date.now() - 42 * 60000).toISOString(),
        end: new Date(Date.now() - 32 * 60000).toISOString()
      }
    },
    // Transcript entries (raw livestream transcript)
    {
      id: 't1',
      timestamp: new Date(Date.now() - 41 * 60000).toISOString(), // 41 minutes ago
      type: 'transcript',
      content: "Let's start by setting up our project with TypeScript. I think it'll give us better type safety as the component library grows.",
      speaker: "Host",
      duration: 12
    },
    {
      id: 't2',
      timestamp: new Date(Date.now() - 40 * 60000).toISOString(), // 40 minutes ago
      type: 'transcript',
      content: "Good idea. What about the build tool? Webpack is traditional but Vite might be faster for development.",
      speaker: "Guest",
      duration: 8
    },
    {
      id: 't3',
      timestamp: new Date(Date.now() - 39 * 60000).toISOString(), // 39 minutes ago
      type: 'transcript',
      content: "I've been really impressed with Vite lately. The HMR is incredibly fast and the configuration is much simpler. Let's go with that.",
      speaker: "Host",
      duration: 15
    },
    {
      id: 't4',
      timestamp: new Date(Date.now() - 38 * 60000).toISOString(), // 38 minutes ago
      type: 'transcript',
      content: "Sounds good. For styling, are we considering CSS modules, styled-components, or something else?",
      speaker: "Guest",
      duration: 7
    },
    {
      id: 't5',
      timestamp: new Date(Date.now() - 21 * 60000).toISOString(), // 21 minutes ago
      type: 'transcript',
      content: "Now that we've decided on the design tokens, let's talk about how we'll structure our typography system. I'm thinking we should define a clear scale.",
      speaker: "Host",
      duration: 14
    },
    {
      id: 't6',
      timestamp: new Date(Date.now() - 20 * 60000).toISOString(), // 20 minutes ago
      type: 'transcript',
      content: "Definitely. A typographic scale with clear ratios will make our design system more consistent. Let's define text sizes, line heights, and letter spacing as part of it.",
      speaker: "Guest",
      duration: 13
    },
    {
      id: '12',
      timestamp: new Date(Date.now() - 12 * 60000).toISOString(), // 12 minutes ago
      type: 'paragraph',
      title: 'Design System Planning',
      content: "We spent considerable time discussing the design token system that will form the foundation of our component library. The conversation centered around creating a scalable and maintainable approach to design variables.\n\nWe explored color systems, typography scales, spacing units, and responsive breakpoints. The team agreed that having a clear naming convention for these tokens is crucial for maintainability.\n\nA decision was made to implement a dark mode from the beginning rather than retrofitting it later, which guided our approach to color variable definition.\n\nThe implementation will use CSS custom properties (variables) to allow runtime theme switching without requiring a rebuild of the application.",
      timeRange: {
        start: new Date(Date.now() - 22 * 60000).toISOString(),
        end: new Date(Date.now() - 14 * 60000).toISOString()
      }
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

// Async thunks for API operations

// Stream info thunks
const fetchStreamInfo = createAsyncThunk(
  'stream/fetchStreamInfo',
  async (_, { dispatch }) => {
    const response = await dispatch(streamApi.endpoints.getStreamInfo.initiate());
    return response.data;
  }
);

const updateStreamInfo = createAsyncThunk(
  'stream/updateStreamInfo',
  async (streamInfo: Partial<StreamInfo>, { dispatch }) => {
    const response = await dispatch(streamApi.endpoints.updateStreamInfo.initiate(streamInfo));
    return response.data;
  }
);

// Steps thunks
const fetchAllSteps = createAsyncThunk(
  'stream/fetchAllSteps',
  async (_, { dispatch }) => {
    const response = await dispatch(stepsApi.endpoints.getAllSteps.initiate());
    return response.data;
  }
);

const addNewUpcomingStep = createAsyncThunk(
  'stream/addNewUpcomingStep',
  async (step: string, { dispatch }) => {
    const response = await dispatch(stepsApi.endpoints.addUpcomingStep.initiate(step));
    return response.data;
  }
);

const setNewActiveStep = createAsyncThunk(
  'stream/setNewActiveStep',
  async (step: string, { dispatch }) => {
    const response = await dispatch(stepsApi.endpoints.setActiveStep.initiate(step));
    return response.data;
  }
);

const completeActiveStep = createAsyncThunk(
  'stream/completeActiveStep',
  async (_, { dispatch }) => {
    const response = await dispatch(stepsApi.endpoints.completeCurrentStep.initiate());
    return response.data;
  }
);

const reactivateStepFromSource = createAsyncThunk(
  'stream/reactivateStepFromSource',
  async (payload: {step: string, source: 'upcoming' | 'completed', stepId: string}, { dispatch }) => {
    const response = await dispatch(stepsApi.endpoints.reactivateStep.initiate(payload));
    return response.data;
  }
);

// Transcript thunks
const fetchTranscript = createAsyncThunk(
  'stream/fetchTranscript',
  async (_, { dispatch }) => {
    const response = await dispatch(transcriptApi.endpoints.getTranscript.initiate());
    return response.data;
  }
);

const addNote = createAsyncThunk(
  'stream/addNote',
  async (content: string, { dispatch }) => {
    const response = await dispatch(transcriptApi.endpoints.addNote.initiate(content));
    return response.data;
  }
);

const addParagraph = createAsyncThunk(
  'stream/addParagraph',
  async (data: {content: string; title?: string; timeRange?: {start: string; end: string}}, { dispatch }) => {
    const response = await dispatch(transcriptApi.endpoints.addParagraph.initiate(data));
    return response.data;
  }
);

// GitHub thunks
const connectToGithub = createAsyncThunk(
  'stream/connectToGithub',
  async (data: {token: string; repoOwner: string; repoName: string}, { dispatch }) => {
    await dispatch(githubApi.endpoints.connectGithub.initiate(data));
    return data;
  }
);

const fetchGithubInfo = createAsyncThunk(
  'stream/fetchGithubInfo',
  async (_, { dispatch }) => {
    const response = await dispatch(githubApi.endpoints.getGithubInfo.initiate());
    return response.data;
  }
);

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
    },
    addTranscriptParagraph: (state, action: PayloadAction<{
      content: string;
      title?: string;
      timeRange?: { start: string; end: string };
    }>) => {
      state.transcript.push({
        id: Math.random().toString(36).substr(2, 9),
        timestamp: new Date().toISOString(),
        type: 'paragraph',
        content: action.payload.content,
        title: action.payload.title,
        timeRange: action.payload.timeRange
      });
    },
    addTranscriptEntry: (state, action: PayloadAction<{
      content: string;
      timestamp: string;
      speaker: string;
      duration?: number;
    }>) => {
      state.transcript.push({
        id: Math.random().toString(36).substr(2, 9),
        timestamp: action.payload.timestamp,
        type: 'transcript',
        content: action.payload.content,
        speaker: action.payload.speaker,
        duration: action.payload.duration
      });
    }
  },
  extraReducers: (builder) => {
    // Add case reducer for resetting state
    builder.addCase('STREAM_RESET_STATE', (state, action: PayloadAction<StreamState>) => {
      // Ensure all parts of the state are reset. Object.assign can be shallow.
      // A more robust way is to return the payload directly if it's the complete initial state.
      return action.payload;
    });

    // Stream info reducers
    builder.addCase(fetchStreamInfo.pending, (state) => {
      state.loading.streamInfo = true;
      state.error.streamInfo = null;
    });
    builder.addCase(fetchStreamInfo.fulfilled, (state, action) => {
      state.loading.streamInfo = false;
      state.info = action.payload;
    });
    builder.addCase(fetchStreamInfo.rejected, (state, action) => {
      state.loading.streamInfo = false;
      state.error.streamInfo = action.error.message || 'Failed to fetch stream info';
    });

    builder.addCase(updateStreamInfo.pending, (state) => {
      state.loading.streamInfo = true;
      state.error.streamInfo = null;
    });
    builder.addCase(updateStreamInfo.fulfilled, (state, action) => {
      state.loading.streamInfo = false;
      state.info = action.payload;
    });
    builder.addCase(updateStreamInfo.rejected, (state, action) => {
      state.loading.streamInfo = false;
      state.error.streamInfo = action.error.message || 'Failed to update stream info';
    });

    // Steps reducers
    builder.addCase(fetchAllSteps.pending, (state) => {
      state.loading.steps = true;
      state.error.steps = null;
    });
    builder.addCase(fetchAllSteps.fulfilled, (state, action) => {
      state.loading.steps = false;
      state.completedSteps = action.payload.completedSteps;
      state.activeStep = action.payload.activeStep;
      state.upcomingSteps = action.payload.upcomingSteps;
      state.stepIdMapping = action.payload._stepIdMapping;
    });
    builder.addCase(fetchAllSteps.rejected, (state, action) => {
      state.loading.steps = false;
      state.error.steps = action.error.message || 'Failed to fetch steps';
    });

    builder.addCase(addNewUpcomingStep.fulfilled, (state, action) => {
      state.upcomingSteps = action.payload.upcomingSteps;
      state.stepIdMapping = action.payload._stepIdMapping;
    });

    builder.addCase(setNewActiveStep.fulfilled, (state, action) => {
      state.completedSteps = action.payload.completedSteps;
      state.activeStep = action.payload.activeStep;
      state.stepIdMapping = action.payload._stepIdMapping;
    });

    builder.addCase(completeActiveStep.fulfilled, (state, action) => {
      state.completedSteps = action.payload.completedSteps;
      state.activeStep = action.payload.activeStep;
      state.upcomingSteps = action.payload.upcomingSteps;
      state.stepIdMapping = action.payload._stepIdMapping;
    });

    builder.addCase(reactivateStepFromSource.fulfilled, (state, action) => {
      state.completedSteps = action.payload.completedSteps;
      state.activeStep = action.payload.activeStep;
      state.upcomingSteps = action.payload.upcomingSteps;
      state.stepIdMapping = action.payload._stepIdMapping;
    });

    // Transcript reducers
    builder.addCase(fetchTranscript.pending, (state) => {
      state.loading.transcript = true;
      state.error.transcript = null;
    });
    builder.addCase(fetchTranscript.fulfilled, (state, action) => {
      state.loading.transcript = false;
      state.transcript = action.payload;
    });
    builder.addCase(fetchTranscript.rejected, (state, action) => {
      state.loading.transcript = false;
      state.error.transcript = action.error.message || 'Failed to fetch transcript';
    });

    builder.addCase(addNote.fulfilled, (state, action) => {
      state.transcript.push(action.payload);
    });

    builder.addCase(addParagraph.fulfilled, (state, action) => {
      state.transcript.push(action.payload);
    });

    // GitHub reducers
    builder.addCase(fetchGithubInfo.pending, (state) => {
      state.loading.github = true;
      state.error.github = null;
    });
    builder.addCase(fetchGithubInfo.fulfilled, (state, action) => {
      state.loading.github = false;
      state.github = action.payload;
    });
    builder.addCase(fetchGithubInfo.rejected, (state, action) => {
      state.loading.github = false;
      state.error.github = action.error.message || 'Failed to fetch GitHub info';
    });

    builder.addCase(connectToGithub.pending, (state) => {
      state.loading.github = true;
      state.error.github = null;
    });
    builder.addCase(connectToGithub.fulfilled, (state, action) => {
      state.loading.github = false;
      state.github.token = action.payload.token;
      state.github.repoOwner = action.payload.repoOwner;
      state.github.repoName = action.payload.repoName;
      state.github.isConnected = true;
    });
    builder.addCase(connectToGithub.rejected, (state, action) => {
      state.loading.github = false;
      state.error.github = action.error.message || 'Failed to connect to GitHub';
    });
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
  // Transcript/Notes actions
  changeTab,
  addTranscriptNote,
  addTranscriptParagraph,
  addTranscriptEntry
} = streamSlice.actions;

// Export async thunks for API operations
export {
  fetchStreamInfo,
  updateStreamInfo,
  fetchAllSteps,
  addNewUpcomingStep,
  setNewActiveStep,
  completeActiveStep,
  reactivateStepFromSource,
  fetchTranscript,
  addNote,
  addParagraph,
  connectToGithub,
  fetchGithubInfo
};

export default streamSlice.reducer;