// Initialize MSW for Storybook only
async function initMocks() {
  // Only used in Storybook environment
  const isStorybook = import.meta.env.STORYBOOK === 'true';
  
  if (isStorybook) {
    console.log('[MSW] Initializing mock service worker for Storybook');
    
    if (typeof window === 'undefined') {
      // Server-side mocking (if needed)
      const { server } = await import('./server');
      server.listen();
    } else {
      // Browser-side mocking
      const { worker } = await import('./browser');
      return worker.start({
        onUnhandledRequest: 'bypass', // Don't warn about unhandled requests
      });
    }
  }
  
  return Promise.resolve();
}

export { initMocks }; 