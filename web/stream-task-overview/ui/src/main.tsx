import React from 'react';
import ReactDOM from 'react-dom/client';
import { Provider } from 'react-redux';
import { store } from './store';
import App from './App';
import './index.css';

// Initialize MSW for development
async function bootstrap() {
  // Enable MSW in development and Storybook
  if (import.meta.env.DEV || import.meta.env.STORYBOOK === 'true') {
    console.log('[App] Initializing MSW in development mode');
    const { initMocks } = await import('./mocks/init');
    await initMocks();
  }

  // Render the application
  ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <Provider store={store}>
        <App />
      </Provider>
    </React.StrictMode>,
  );
}

// Start the application
bootstrap().catch(err => console.error('[App] Failed to start application:', err));