import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      // Proxy all /auth requests to the backend server
      '/auth': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      },
      // You may also want to proxy API requests
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        secure: false
      }
    }
  }
})