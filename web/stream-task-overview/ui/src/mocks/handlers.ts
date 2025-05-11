import { http, HttpResponse } from 'msw';

export const handlers = [
  // Mock authentication status
  http.get('/api/auth/status', () => {
    return HttpResponse.json({ 
      isAuthenticated: true, 
      isAdmin: true 
    });
  }),

  // Get user information
  http.get('/api/auth/user', () => {
    return HttpResponse.json({
      authenticated: true,
      id: 'user123',
      login: 'admin',
      name: 'Admin User',
      email: 'admin@example.com',
      avatar: 'https://avatar.example.com/admin',
      is_admin: true
    });
  }),
  
  // Mock steps data
  http.get('/api/stream/steps', () => {
    return HttpResponse.json({
      completedSteps: ['Research competitors', 'Create wireframes'],
      activeStep: 'Implement UI components',
      upcomingSteps: ['Write unit tests', 'Deploy to staging'],
      stepIdMapping: {
        'Research competitors': '1',
        'Create wireframes': '2',
        'Implement UI components': '3',
        'Write unit tests': '4',
        'Deploy to staging': '5'
      }
    });
  }),
  
  // Mock stream info endpoint for streamApi
  http.get('/api/stream', () => {
    return HttpResponse.json({
      title: 'Building a Task Management System',
      description: 'Creating a full-stack application with React and Go',
      language: 'TypeScript/Go',
      githubRepo: 'organization/repo-name',
      startTime: '2025-05-11T18:30:00.000Z',
      viewerCount: 128
    });
  }),
  
  // Mock update stream info for streamApi
  http.put('/api/stream', async () => {
    return HttpResponse.json({
      title: 'Building a Task Management System',
      description: 'Creating a full-stack application with React and Go',
      language: 'TypeScript/Go',
      githubRepo: 'organization/repo-name',
      startTime: '2025-05-11T18:30:00.000Z',
      viewerCount: 128
    });
  }),

  // GitHub related endpoints
  http.get('/api/github/info', () => {
    return HttpResponse.json({
      repoUrl: 'https://github.com/organization/repo-name',
      isConnected: true,
      token: 'mock-github-token',
      repoOwner: 'organization',
      repoName: 'repo-name',
      currentBranch: 'main',
      latestCommit: {
        message: 'Initial commit',
        author: 'Developer',
        hash: 'abc123',
        date: '2023-05-01T10:00:00.000Z',
        url: 'https://github.com/organization/repo-name/commit/abc123'
      }
    });
  }),

  http.get('/api/github/commits', () => {
    return HttpResponse.json([
      {
        message: 'Initial commit',
        author: 'Developer',
        hash: 'abc123',
        date: '2023-05-01T10:00:00.000Z',
        url: 'https://github.com/organization/repo-name/commit/abc123'
      }
    ]);
  }),
  
  // Transcript endpoints
  http.get('/api/stream/transcript', () => {
    return HttpResponse.json([
      {
        type: 'system',
        content: 'Stream started',
        timestamp: '2023-05-01T10:00:00.000Z'
      },
      {
        type: 'note',
        content: 'This is a test note',
        timestamp: '2023-05-01T10:05:00.000Z'
      }
    ]);
  }),

  // Mock login endpoint
  http.post('/api/auth/login', async ({ request }) => {
    const body = await request.json();
    if (body.username === 'admin' && body.password === 'password') {
      return HttpResponse.json({ 
        success: true, 
        isAdmin: true,
        token: 'mock-jwt-token' 
      });
    }
    return new HttpResponse(
      JSON.stringify({ success: false, message: 'Invalid credentials' }),
      { status: 401 }
    );
  }),

  // Mock logout endpoint
  http.post('/api/auth/logout', () => {
    return HttpResponse.json({ success: true });
  })
];