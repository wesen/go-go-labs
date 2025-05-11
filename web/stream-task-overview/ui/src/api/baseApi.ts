import { fetchBaseQuery } from '@reduxjs/toolkit/query/react';

// Base API configuration
export const baseQuery = fetchBaseQuery({
  baseUrl: 'http://localhost:8080/api',
  prepareHeaders: (headers) => {
    // Add common headers here if needed
    return headers;
  },
});

// Common error transformer
export const transformErrorResponse = (response: any, meta: any, arg: any) => {
  return {
    status: response.status,
    message: response.data?.error || 'An unexpected error occurred',
  };
};

// Utility functions for data transformation

// Convert snake_case to camelCase (backend to frontend)
export const toCamelCase = (obj: any): any => {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(item => toCamelCase(item));
  }

  return Object.keys(obj).reduce((acc, key) => {
    const camelKey = key.replace(/(_\w)/g, m => m[1].toUpperCase());
    acc[camelKey] = toCamelCase(obj[key]);
    return acc;
  }, {} as any);
};

// Convert camelCase to snake_case (frontend to backend)
export const toSnakeCase = (obj: any): any => {
  if (obj === null || typeof obj !== 'object') {
    return obj;
  }

  if (Array.isArray(obj)) {
    return obj.map(item => toSnakeCase(item));
  }

  return Object.keys(obj).reduce((acc, key) => {
    const snakeKey = key.replace(/([A-Z])/g, m => `_${m.toLowerCase()}`);
    acc[snakeKey] = toSnakeCase(obj[key]);
    return acc;
  }, {} as any);
};