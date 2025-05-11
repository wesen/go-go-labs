# GitHub OAuth and Polling Implementation

## Overview

This document explains the technical implementation of GitHub OAuth authentication and GitHub API polling in the Stream Task Viewer application.

## Authentication Implementation

The authentication system is built using:

- **labstack/echo-contrib/session**: For cookie-based session management
- **markbates/goth**: For OAuth provider integration
- **agentstation/egothic**: For adapting Goth to Echo framework
- **gorilla/sessions**: For session store implementation

### Key Components

1. **auth.go**:
   - Manages GitHub OAuth flow 
   - Provides admin authorization 
   - Handles session management
   - Implements middleware for access control

2. **Middleware**:
   - `RequireAuth`: Ensures authenticated access to protected endpoints
   - `AdminOnly`: Restricts routes to admin users only

3. **Sessions**:
   - Cookie-based sessions with encryption
   - Stores GitHub user details and admin status
   - HTTP-only cookies with secure flags

## File Structure

```
backend/
  ├── auth.go            - Authentication implementation
  ├── github_poller.go   - GitHub API polling implementation
  ├── main.go            - Server setup and route configuration
  └── docs/
      ├── github-oauth-setup.md    - Setup instructions
      └── oauth-implementation.md  - Implementation details
```

## Admin Authentication

Admin access is determined by comparing the authenticated GitHub username with the `ADMIN_GITHUB_USERNAME` environment variable. When a user logs in, the system:

1. Retrieves their GitHub username (NickName)
2. Compares it case-insensitively with the configured admin username
3. Sets an `is_admin` flag in the session if it matches

This approach allows controlled admin access without modifying the database schema.

```go
// Example code from auth.go
if h.adminUsername != "" && strings.EqualFold(user.NickName, h.adminUsername) {
    log.Info().Str("username", user.NickName).Msg("Admin user authenticated")
    sess.Values["is_admin"] = true
}
```

## GitHub Viewer Count Polling

The `GithubPoller` uses the GitHub API to periodically fetch repository viewer counts and update the stream information.

### Polling Implementation

1. **Initialization**:
   - Creates a background goroutine for polling
   - Configurable interval (default: 2 minutes)
   - Uses `time.Ticker` for regular updates

2. **API Integration**:
   - Uses the GitHub Traffic API (`/repos/{owner}/{repo}/traffic/views`)
   - Requires a personal access token with proper permissions
   - Extracts repository path from the stream info's GitHub URL

3. **Data Flow**:
   - Fetches view count data from GitHub
   - Updates the StreamInfo.ViewerCount field
   - Persists changes to the database

### Technical Considerations

- **Rate Limiting**: Default 2-minute interval helps avoid GitHub API rate limits
- **Error Handling**: Robust error handling for network issues
- **Thread Safety**: Uses mutex for concurrent access protection
- **Graceful Shutdown**: Properly terminates polling on application shutdown
- **URL Parsing**: Handles various GitHub URL formats

## API Routes

The updated API structure includes:

### Public Routes
- `GET /api/stream`: View stream information
- `GET /api/stream/steps`: View all steps
- `GET /auth/user`: Get current user information

### Protected Routes (Admin Only)
- `PUT /api/stream`: Update stream information
- `PUT /api/stream/steps/active`: Set active step
- `POST /api/stream/steps/upcoming`: Add upcoming step
- `POST /api/stream/steps/complete`: Complete active step
- `PUT /api/stream/steps/reactivate`: Reactivate a step

### Authentication Routes
- `GET /auth/github`: Begin GitHub OAuth flow
- `GET /auth/github/callback`: Handle OAuth callback
- `GET /auth/logout`: Log out current user

## Frontend Integration

The React frontend communicates with the authentication system via the `/auth/user` endpoint and renders different UI components based on authentication status and admin privileges.

The `AuthStatus` component:
- Displays login state
- Shows GitHub user info when authenticated
- Indicates admin status with a shield icon
- Provides login/logout functionality
- Displays user's avatar when available

## Session Data Structure

User sessions contain the following information:

```json
{
  "user_id": "12345678",
  "login": "github_username",
  "name": "User's Full Name",
  "email": "user@example.com",
  "avatar": "https://github.com/user.png",
  "is_admin": true
}
```

## Future Improvements

- Implement additional OAuth providers (Google, GitLab, etc.)
- Add role-based access control for more granular permissions
- Implement rate limiting for API endpoints
- Add JWT token support for mobile/SPA clients
- Enhance security with CSRF protection and strict CSP
- Add more detailed logging and monitoring