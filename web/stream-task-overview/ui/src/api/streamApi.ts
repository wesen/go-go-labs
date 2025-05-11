import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQuery, toCamelCase, toSnakeCase } from './baseApi';
import { StreamInfo } from '../store/slices/streamSlice';

// Transform helpers for StreamInfo
export const toBackendStreamInfo = (frontendInfo: Partial<StreamInfo>) => {
  return toSnakeCase({
    title: frontendInfo.title,
    description: frontendInfo.description,
    startTime: frontendInfo.startTime,
    language: frontendInfo.language,
    githubRepo: frontendInfo.githubRepo,
    viewerCount: frontendInfo.viewerCount
  });
};

export const toFrontendStreamInfo = (backendInfo: any): StreamInfo => {
  const camelCased = toCamelCase(backendInfo);
  return {
    title: camelCased.title,
    description: camelCased.description,
    startTime: camelCased.startTime,
    language: camelCased.language,
    githubRepo: camelCased.githubRepo,
    currentTask: camelCased.activeStep || '',
    viewerCount: camelCased.viewerCount
  } as StreamInfo;
};

export const streamApi = createApi({
  reducerPath: 'streamApi',
  baseQuery,
  endpoints: (builder) => ({
    getStreamInfo: builder.query<StreamInfo, void>({
      query: () => 'stream',
      transformResponse: (response) => toFrontendStreamInfo(response)
    }),
    updateStreamInfo: builder.mutation<StreamInfo, Partial<StreamInfo>>({
      query: (streamInfo) => ({
        url: 'stream',
        method: 'PUT',
        body: toBackendStreamInfo(streamInfo),
      }),
      transformResponse: (response) => toFrontendStreamInfo(response)
    }),
  }),
});

export const { 
  useGetStreamInfoQuery, 
  useUpdateStreamInfoMutation 
} = streamApi;