# Stream Task Overview

An application for tracking and displaying programming stream tasks, progress and GitHub integration.

## Features

- Stream information display
- Step tracking (active, upcoming, completed)
- Transcript entries with various types
- GitHub integration
- Admin functionality protected by GitHub OAuth

## Environment Variables

To run this application, set the following environment variables:

```bash
# GitHub OAuth (required for auth)
GITHUB_CLIENT_ID=your_github_client_id
GITHUB_CLIENT_SECRET=your_github_client_secret
GITHUB_CALLBACK_URL=http://localhost:8080/auth/github/callback

# Admin access
ADMIN_GITHUB_USERNAME=your_github_username

# Session security
SESSION_SECRET=your_session_secret

# Server configuration
PORT=8080
ENV=development
```

## Development Setup

### Backend (Go)

```bash
cd backend
go mod download
go run main.go
```

### Frontend (React)

```bash
cd ui
bun install
bun run dev
```

## Admin Features

Admin features are protected behind GitHub authentication. Only users with the GitHub username specified in `ADMIN_GITHUB_USERNAME` will have access to admin features including:

- Adding and managing tasks
- Adding transcript notes
- Updating stream information
- Connecting to GitHub repositories