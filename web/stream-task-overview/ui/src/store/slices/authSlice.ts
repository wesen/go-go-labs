import { createSlice, PayloadAction, SerializedError, createAction } from '@reduxjs/toolkit';
import { authApi } from '../../api/authApi';

interface AuthState {
  isAuthenticated: boolean;
  isAdmin: boolean;
  token: string | null;
  loading: boolean;
  error: string | null;
}

export const initialState: AuthState = {
  isAuthenticated: false,
  isAdmin: false,
  token: typeof localStorage !== 'undefined' ? localStorage.getItem('token') : null,
  loading: false,
  error: null,
};

// Define a type for the login error payload if it's a FetchBaseQueryError
interface LoginErrorPayload {
  status: number;
  data: {
    message?: string;
    // other properties the backend might send on error
  };
}

// Create an action creator for resetting the auth state
export const resetAuthState = createAction<AuthState>('AUTH_RESET_STATE');

const authSlice = createSlice({
  name: 'auth',
  initialState,
  reducers: {
    setCredentials: (state, action: PayloadAction<{ token: string; isAdmin: boolean }>) => {
      const { token, isAdmin } = action.payload;
      state.token = token;
      state.isAuthenticated = true;
      state.isAdmin = isAdmin;
      localStorage.setItem('token', token);
    },
    logout: (state) => {
      state.token = null;
      state.isAuthenticated = false;
      state.isAdmin = false;
      if (typeof localStorage !== 'undefined') {
        localStorage.removeItem('token');
      }
    },
  },
  extraReducers: (builder) => {
    // Add case reducer for resetting state using the action creator
    builder.addCase(resetAuthState, (state, action) => {
      // action is now correctly PayloadAction<AuthState>
      return action.payload;
    });

    // Add reducers for RTK Query API endpoints
    builder
      // Handle login success
      .addMatcher(
        authApi.endpoints.login.matchFulfilled,
        (state, { payload }) => {
          if (payload.success && payload.token) {
            state.token = payload.token;
            state.isAuthenticated = true;
            state.isAdmin = payload.isAdmin;
            state.loading = false;
            state.error = null;
            // Store token in localStorage for persistence
            localStorage.setItem('token', payload.token);
          }
        }
      )
      // Handle login failure
      .addMatcher(
        authApi.endpoints.login.matchRejected,
        (state, action) => {
          state.loading = false;
          if (action.payload) { // payload can be undefined if the error is not from fetchBaseQuery
            if ('data' in action.payload && action.payload.data && typeof action.payload.data === 'object') {
              // This is likely a FetchBaseQueryError
              const errorData = action.payload.data as { message?: string }; 
              state.error = errorData.message || 'Failed to log in';
            } else if ('message' in action.payload) {
              // This could be a SerializedError or other error with a message property
              state.error = (action.payload as SerializedError).message || 'Failed to log in';
            } else {
              state.error = 'An unknown error occurred during login.';
            }
          } else if (action.error && action.error.message) {
            // Fallback to action.error if payload is not informative (e.g. network error)
            state.error = action.error.message;
          } else {
            state.error = 'Failed to log in due to an unknown error.';
          }
        }
      )
      // Handle logout success
      .addMatcher(
        authApi.endpoints.logout.matchFulfilled,
        (state) => {
          state.token = null;
          state.isAuthenticated = false;
          state.isAdmin = false;
          if (typeof localStorage !== 'undefined') { // Added SSR check here too for consistency
            localStorage.removeItem('token');
          }
        }
      )
      // Handle auth status check success
      .addMatcher(
        authApi.endpoints.getAuthStatus.matchFulfilled,
        (state, { payload }) => {
          state.isAuthenticated = payload.isAuthenticated;
          state.isAdmin = payload.isAdmin;
          state.loading = false;
        }
      );
  },
});

export const { setCredentials, logout } = authSlice.actions;
export default authSlice.reducer;