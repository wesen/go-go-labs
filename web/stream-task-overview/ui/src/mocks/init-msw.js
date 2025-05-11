// Initialize MSW
// Run this script once to set up the worker file

import { worker } from './browser';

// Initialize the browser worker
export async function initMsw() {
  console.log('[MSW] Initializing MSW browser worker');
  
  if (process.env.NODE_ENV === 'development' || process.env.STORYBOOK === 'true') {
    return worker.start({
      onUnhandledRequest: 'bypass', // Don't warn on unhandled requests
    });
  }
  
  return Promise.resolve();
} 