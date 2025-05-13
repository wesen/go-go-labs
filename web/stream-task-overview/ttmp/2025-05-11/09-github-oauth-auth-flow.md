# GitHub OAuth Authentication Flow

This document explains how GitHub OAuth authentication works in the application, detailing both the backend and frontend implementations.

## Backend Implementation (Go)

The backend uses the following libraries for GitHub OAuth:
- `github.com/markbates/goth` - A library for OAuth authentication
- `github.com/agentstation/egothic` - An Echo-compatible adapter for Goth
- `github.com/gorilla/sessions` - For session management
- `github.com/labstack/echo-contrib/session` - Echo middleware for sessions

### Configuration

The OAuth configuration requires three environment variables:
- `GITHUB_CLIENT_ID` - Your GitHub OAuth application client ID
- `GITHUB_CLIENT_SECRET` - Your GitHub OAuth application client secret
- `GITHUB_CALLBACK_URL` - The callback URL (defaults to http://localhost:8080/auth/github/callback)

Additionally, an optional `ADMIN_GITHUB_USERNAME` can be set to designate a GitHub username as an admin.

### Authentication Flow

1. **Initialization**: The OAuth setup is initialized in `main.go`:
   ```go
   InitAuth(e)
   auth := NewAuthHandler()
   ```

2. **Route Setup**: Four main auth routes are registered:
   ```go
   e.GET("/auth/github", auth.GithubAuthBegin)
   e.GET("/auth/github/callback", auth.GithubAuthCallback)
   e.GET("/auth/user", auth.GetCurrentUser)
   e.GET("/auth/logout", auth.Logout)
   ```

3. **OAuth Flow**:
   - **Start Auth** (`/auth/github`): When a user clicks "Login with GitHub", they're redirected to this endpoint, which initiates the GitHub OAuth flow
   - **Callback** (`/auth/github/callback`): GitHub redirects the user here after authorization with a code
   - **Session Storage**: Upon successful authentication, user info is stored in the session:
     ```go
     sess.Values["user_id"] = user.UserID
     sess.Values["login"] = user.NickName
     sess.Values["name"] = user.Name
     sess.Values["email"] = user.Email
     sess.Values["avatar"] = user.AvatarURL
     ```
   - **Admin Check**: If the authenticated user's GitHub username matches `ADMIN_GITHUB_USERNAME`, they get admin privileges:
     ```go
     if h.adminUsername != "" && strings.EqualFold(user.NickName, h.adminUsername) {
         sess.Values["is_admin"] = true
     }
     ```

4. **Session Management**:
   - Sessions are stored as secure cookies with a 7-day expiration
   - Session secret is configurable via `SESSION_SECRET` env var

5. **Authorization Scopes**:
   The application requests the following GitHub OAuth scopes:
   - `read:user` - Read access to user profile data
   - `user:email` - Access to user email addresses

## Frontend Implementation (React)

The frontend uses React with Redux Toolkit to manage the OAuth flow.

### Components

1. **LoginForm Component**: 
   - The `LoginForm.tsx` component is a traditional username/password form
   - This is NOT used for GitHub OAuth authentication
   - GitHub OAuth uses a separate mechanism

2. **Authentication API**:
   - The app defines an `authApi` using RTK Query in `authApi.ts`
   - It includes endpoints for checking auth status
   - The base URL is configured in `baseApi.ts` (default: http://localhost:8080/api)

3. **Auth Custom Hook**:
   - `useAuth` hook in `useAuth.tsx` provides authentication functions
   - Wraps RTK Query endpoints
   - Exposes `isAuthenticated` and `isAdmin` states

### OAuth Flow in Frontend

Since the GitHub OAuth flow is handled at the browser level with redirects:

1. **Login Initiation**:
   - A "Login with GitHub" button in the UI should link directly to `/auth/github`
   - No AJAX call is needed as this is a full page redirect
   - Example:
     ```jsx
     <a href="/auth/github" className="github-login-btn">Login with GitHub</a>
     ```

2. **After Authentication**:
   - After auth, the user is redirected back to the frontend (typically the homepage)
   - The app should check authentication status by making a GET request to `/auth/user`
   - This is likely handled in the app's main layout component on mount

### Important Note on Auth Endpoints

There's a discrepancy between the backend and frontend code regarding authentication endpoints:

1. **Backend Endpoints**:
   - `/auth/github` - Initiates GitHub OAuth flow (implemented)
   - `/auth/github/callback` - Handles GitHub OAuth callback (implemented)
   - `/auth/user` - Gets current user information (implemented)
   - `/auth/logout` - Logs out the user (implemented)

2. **Frontend References**:
   - The `AuthStatus.tsx` component contains a login link to `/auth/login` instead of `/auth/github`
   - The `authApi.ts` contains endpoints for `auth/login`, which is for username/password login

3. **Resolution**:
   - To use GitHub OAuth, modify the login button in `AuthStatus.tsx` to point to `/auth/github` instead of `/auth/login`:
     ```jsx
     <a 
       href="/auth/github" 
       className="px-3 py-1 bg-black text-white text-xs uppercase tracking-wider inline-flex items-center"
     >
       <User size={12} className="mr-1" /> Login with GitHub
     </a>
     ```
   - The username/password login form (`LoginForm.tsx`) and its API endpoint (`auth/login`) appear to be placeholder code and not fully implemented in the backend

## Setting Up GitHub OAuth

To set up GitHub OAuth for your application:

1. **Create a GitHub OAuth App**:
   - Go to GitHub Settings > Developer Settings > OAuth Apps > New OAuth App
   - Fill in:
     - **Application Name**: Your app name
     - **Homepage URL**: Your app URL (e.g., http://localhost:8080)
     - **Authorization callback URL**: Your callback URL (e.g., http://localhost:8080/auth/github/callback)

2. **Set Environment Variables**:
   ```bash
   export GITHUB_CLIENT_ID="your_client_id"
   export GITHUB_CLIENT_SECRET="your_client_secret"
   export GITHUB_CALLBACK_URL="http://localhost:8080/auth/github/callback"  # Optional if using default
   export ADMIN_GITHUB_USERNAME="your_github_username"  # Optional for admin access
   export SESSION_SECRET="random_secret_string"  # Optional but recommended for security
   ```

3. **Configure Frontend**:
   - Ensure the frontend has a "Login with GitHub" button that links to `/auth/github`
   - After login, use the `useGetAuthStatusQuery` hook to get user information

## Initializing GitHub Integration in the Database

The error message `Failed to get GitHub info from database error="sql: no rows in result set"` occurs because there is no GitHub integration data in your database yet. To fix this, you need to connect your application to GitHub:

### Method 1: Using the Admin API

The easiest way to initialize GitHub integration is by using the API endpoint designed for this purpose:

1. First, ensure you're logged in as an admin (your GitHub username matches `ADMIN_GITHUB_USERNAME`)

2. Make a POST request to `/api/github/connect` with the following JSON payload:
   ```json
   {
     "token": "your_github_personal_access_token",
     "repoOwner": "username_or_organization",
     "repoName": "repository_name"
   }
   ```

   You can use cURL or any API client:
   ```bash
   curl -X POST http://localhost:8080/api/github/connect \
     -H "Content-Type: application/json" \
     -d '{"token":"ghp_your_personal_access_token","repoOwner":"wesen","repoName":"your-repo-name"}'
   ```

3. The backend will store this information in the `github_integration` table, which will resolve the error

### Method 2: Direct SQL Insertion

If you can't access the API (e.g., because you're not authenticated as admin yet), you can insert the data directly into the SQLite database:

1. Open the SQLite database:
   ```bash
   sqlite3 backend/stream.db
   ```

2. Insert a record into the `github_integration` table:
   ```sql
   INSERT INTO github_integration (token, repo_owner, repo_name, current_branch)
   VALUES ('your_github_personal_access_token', 'wesen', 'your-repo-name', 'main');
   ```

3. Exit SQLite:
   ```
   .exit
   ```

### GitHub Personal Access Token

For either method, you'll need a GitHub Personal Access Token with the appropriate permissions:

1. Go to GitHub Settings > Developer Settings > Personal Access Tokens > Tokens (classic)
2. Click "Generate new token" and select "Generate new token (classic)"
3. Give it a descriptive name
4. Select the following scopes:
   - `repo` (Full control of private repositories)
   - `read:user` (Read access to user profile data)
5. Click "Generate token" and copy the token immediately (you won't be able to see it again)

### Database Schema

The GitHub integration information is stored in the `github_integration` table with the following structure:

```sql
CREATE TABLE IF NOT EXISTS github_integration (
  id INTEGER PRIMARY KEY,
  token TEXT NOT NULL,
  repo_owner TEXT NOT NULL,
  repo_name TEXT NOT NULL,
  current_branch TEXT NOT NULL,
  latest_commit_hash TEXT,
  latest_commit_message TEXT,
  latest_commit_author TEXT,
  latest_commit_date DATETIME,
  latest_commit_url TEXT
);
```

Once you've initialized the GitHub integration data, the application will be able to retrieve GitHub information and the error should be resolved.

## Authentication Flow Diagram

```
┌──────────┐    1. Click GitHub Login     ┌──────────┐
│          │ ─────────────────────────────▶          │
│  Browser │                               │ Backend  │
│          │ ◀─────────────────────────────│          │
└──────────┘    2. Redirect to GitHub      └──────────┘
     │                                          ▲
     │                                          │
     │                                          │
     ▼                                          │
┌──────────┐    3. Authorize App           │
│          │                               │
│  GitHub  │                               │
│          │    4. Redirect to callback    │
└──────────┘ ─────────────────────────────▶
```

## Security Considerations

1. Keep `GITHUB_CLIENT_SECRET` secure and never expose it in frontend code
2. Use HTTPS in production to protect cookies and OAuth data
3. Consider additional security measures like rate limiting on auth endpoints

## Troubleshooting

- **GitHub OAuth not working**:
  - Check environment variables are correctly set
  - Verify callback URL matches the one registered in GitHub OAuth app settings
  - Check browser console and server logs for errors

- **Admin access not working**:
  - Verify `ADMIN_GITHUB_USERNAME` is set correctly
  - Note that the username comparison is case-insensitive but must otherwise match exactly 

- **"Failed to get GitHub info from database" error**:
  - This error occurs when no GitHub integration data exists in the database
  - Use the instructions in the "Initializing GitHub Integration" section above to resolve it 