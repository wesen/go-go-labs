import { createApi } from '@reduxjs/toolkit/query/react';
import { baseQuery, toCamelCase, toSnakeCase } from './baseApi';
import { TranscriptEntry } from '../store/slices/streamSlice';

export const transcriptApi = createApi({
  reducerPath: 'transcriptApi',
  baseQuery,
  tagTypes: ['Transcript'],
  endpoints: (builder) => ({
    getTranscript: builder.query<TranscriptEntry[], void>({
      query: () => 'stream/transcript',
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry[],
      providesTags: ['Transcript'],
    }),
    getTranscriptByType: builder.query<TranscriptEntry[], string>({
      query: (type) => `stream/transcript/types/${type}`,
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry[],
      providesTags: ['Transcript'],
    }),
    addNote: builder.mutation<TranscriptEntry, string>({
      query: (content) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'note',
          content,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
    addParagraph: builder.mutation<TranscriptEntry, {
      content: string;
      title?: string;
      timeRange?: { start: string; end: string };
    }>({
      query: (data) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'paragraph',
          ...data,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
    addTaskStarted: builder.mutation<TranscriptEntry, { taskName: string }>({
      query: ({ taskName }) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'task_started',
          content: `Started working on ${taskName}`,
          taskName,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
    addTaskCompleted: builder.mutation<TranscriptEntry, { taskName: string }>({
      query: ({ taskName }) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'task_completed',
          content: `Completed ${taskName}`,
          taskName,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
    addCommit: builder.mutation<TranscriptEntry, { 
      content: string;
      commitHash: string;
      commitUrl: string;
    }>({
      query: (data) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'commit',
          ...data,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
    addTranscriptEntry: builder.mutation<TranscriptEntry, {
      content: string;
      speaker: string;
      duration?: number;
    }>({
      query: (data) => ({
        url: 'stream/transcript',
        method: 'POST',
        body: toSnakeCase({ 
          type: 'transcript',
          ...data,
          timestamp: new Date().toISOString() 
        }),
      }),
      transformResponse: (response) => toCamelCase(response) as TranscriptEntry,
      invalidatesTags: ['Transcript'],
    }),
  }),
});

export const {
  useGetTranscriptQuery,
  useGetTranscriptByTypeQuery,
  useAddNoteMutation,
  useAddParagraphMutation,
  useAddTaskStartedMutation,
  useAddTaskCompletedMutation,
  useAddCommitMutation,
  useAddTranscriptEntryMutation
} = transcriptApi;