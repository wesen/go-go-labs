// Initialize MSW
async function initMocks() {
  // Check if we should enable mocking
  const shouldMock = import.meta.env.DEV || import.meta.env.STORYBOOK;
  
  if (shouldMock) {
    console.log('[MSW] Initializing mock service worker');
    
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