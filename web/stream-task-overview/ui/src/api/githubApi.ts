import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQuery, toCamelCase, toSnakeCase } from './baseApi';
import { GithubInfo } from '../store/slices/streamSlice';

export const githubApi = createApi({
  reducerPath: 'githubApi',
  baseQuery,
  tagTypes: ['GitHub'],
  endpoints: (builder) => ({
    connectGithub: builder.mutation<void, {
      token: string;
      repoOwner: string;
      repoName: string;
    }>({
      query: (data) => ({
        url: 'github/connect',
        method: 'POST',
        body: toSnakeCase(data),
      }),
      invalidatesTags: ['GitHub'],
    }),
    getGithubInfo: builder.query<GithubInfo, void>({
      query: () => 'github/info',
      transformResponse: (response) => toCamelCase(response) as GithubInfo,
      providesTags: ['GitHub'],
    }),
    getCommits: builder.query<any[], { limit?: number }>({
      query: ({ limit = 10 }) => `github/commits?limit=${limit}`,
      transformResponse: (response) => toCamelCase(response),
      providesTags: ['GitHub'],
    }),
  }),
});

export const {
  useConnectGithubMutation,
  useGetGithubInfoQuery,
  useGetCommitsQuery
} = githubApi;