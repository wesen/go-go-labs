# GitHub OAuth and Authentication Setup

## Overview

This document explains how to set up GitHub OAuth authentication for the Stream Task Viewer application, including admin access and GitHub viewer count polling.

## Required Environment Variables

The following environment variables need to be set for GitHub OAuth to work properly:

- `GITHUB_CLIENT_ID`: OAuth application client ID from GitHub
- `GITHUB_CLIENT_SECRET`: OAuth application client secret from GitHub
- `GITHUB_CALLBACK_URL`: The URL GitHub will redirect to after authentication (defaults to `http://localhost:8080/auth/github/callback`)
- `SESSION_SECRET`: A secret key used to encrypt session cookies (32+ random bytes recommended)
- `ADMIN_GITHUB_USERNAME`: The GitHub username of the user who should have admin privileges
- `GITHUB_API_TOKEN`: A personal access token used for polling repository viewer counts
- `GITHUB_POLLING_INTERVAL`: How often to poll GitHub for viewer counts (optional, defaults to 2 minutes)

## Setting Up a GitHub OAuth App

1. Go to GitHub.com and log in
2. Navigate to Settings → Developer settings → OAuth Apps → "New OAuth App"
3. Fill in the application details:
   - Application name: Stream Task Viewer
   - Homepage URL: http://localhost:8080
   - Authorization callback URL: http://localhost:8080/auth/github/callback
   - Description: (Optional) Stream task viewer for live coding sessions
4. Click "Register application"
5. Copy the Client ID and generate a new Client Secret
6. Set the needed permissions: `read:user` and `user:email` scopes are required

## Setting Up Environment Variables

Create a `.env` file in the root directory with the following contents:

```
GITHUB_CLIENT_ID=your_client_id
GITHUB_CLIENT_SECRET=your_client_secret
ADMIN_GITHUB_USERNAME=your_github_username
SESSION_SECRET=$(openssl rand -base64 32)
GITHUB_API_TOKEN=your_personal_access_token
```

Or export them in your terminal:

```bash
export GITHUB_CLIENT_ID=your_client_id
export GITHUB_CLIENT_SECRET=your_client_secret
export ADMIN_GITHUB_USERNAME=your_github_username
export SESSION_SECRET=$(openssl rand -base64 32)
export GITHUB_API_TOKEN=your_personal_access_token
```

## Authentication Flow

1. User visits `/auth/github` to initiate OAuth flow
2. User is redirected to GitHub's authorization page
3. Upon approval, GitHub redirects back to `/auth/github/callback`
4. The server verifies the token and creates a session
5. The user is redirected to the main application

## Admin Access

Only the GitHub user whose username matches the `ADMIN_GITHUB_USERNAME` environment variable will be granted admin privileges. Admin users can:

- Update stream information
- Add, update, and manage steps
- Access admin-only API endpoints

## GitHub Viewer Count Polling

The application automatically polls the GitHub API to get repository view counts. This requires:

1. A GitHub personal access token with `repo` scope
2. Admin access to the repository specified in the stream information

The polling interval defaults to 2 minutes but can be configured with the `GITHUB_POLLING_INTERVAL` environment variable (e.g., `10m` for 10 minutes).

## Creating a GitHub Personal Access Token

To create a GitHub Personal Access Token for the viewer count polling:

1. Go to GitHub.com and log in
2. Navigate to Settings → Developer settings → Personal access tokens → Tokens (classic)
3. Click "Generate new token" → "Generate new token (classic)"
4. Give your token a descriptive name
5. Select the `repo` scope (specifically `repo:status` and `repo_deployment` are needed for traffic data)
6. Click "Generate token"
7. Copy the token and save it as your `GITHUB_API_TOKEN`

## Security Considerations

- Always use HTTPS in production
- Generate a strong random value for `SESSION_SECRET`
- Keep your OAuth client secret and personal access token secure
- Restrict the scopes of your personal access token to only what's needed
- Consider implementing additional security measures for production deployments
- Use environment variables rather than hardcoding secrets
- Rotate tokens and secrets regularly