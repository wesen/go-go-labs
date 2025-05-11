import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { baseUrl } from './baseApi';

interface AuthResponse {
  success: boolean;
  isAdmin: boolean;
  token?: string;
  message?: string;
}

interface AuthStatus {
  isAuthenticated: boolean;
  isAdmin: boolean;
}

interface LoginRequest {
  username: string;
  password: string;
}

// Create an API slice for authentication
export const authApi = createApi({
  reducerPath: 'authApi',
  baseQuery: fetchBaseQuery({ baseUrl }),
  tagTypes: ['Auth'],
  endpoints: (builder) => ({
    login: builder.mutation<AuthResponse, LoginRequest>({
      query: (credentials) => ({
        url: 'auth/login',
        method: 'POST',
        body: credentials,
      }),
      invalidatesTags: ['Auth'],
    }),
    logout: builder.mutation<void, void>({
      query: () => ({
        url: 'auth/logout',
        method: 'POST',
      }),
      invalidatesTags: ['Auth'],
    }),
    getAuthStatus: builder.query<AuthStatus, void>({
      query: () => 'auth/status',
      providesTags: ['Auth'],
    }),
  }),
});

// Export hooks for using the API
export const {
  useLoginMutation,
  useLogoutMutation,
  useGetAuthStatusQuery,
} = authApi;