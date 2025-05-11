/**
 * GitHub API Client Service
 * Handles authentication and communication with GitHub API v3
 */

// GitHub API base URL
const GITHUB_API_BASE = 'https://api.github.com';

// Types for GitHub API responses
export interface GithubRepository {
  name: string;
  full_name: string;
  html_url: string;
  description: string | null;
  default_branch: string;
  language: string | null;
  stargazers_count: number;
  forks_count: number;
  open_issues_count: number;
}

export interface GithubBranch {
  name: string;
  commit: {
    sha: string;
    url: string;
  };
}

export interface GithubCommit {
  sha: string;
  commit: {
    author: {
      name: string;
      email: string;
      date: string;
    };
    message: string;
  };
  html_url: string;
  author: {
    login: string;
    avatar_url: string;
    html_url: string;
  } | null;
}

export interface GithubApiError {
  message: string;
  documentation_url?: string;
}

// API Client class
class GithubApiClient {
  private token: string | null = null;

  /**
   * Set the GitHub personal access token for authentication
   */
  setToken(token: string) {
    this.token = token;
  }

  /**
   * Clear the authentication token
   */
  clearToken() {
    this.token = null;
  }

  /**
   * Check if token is set
   */
  hasToken(): boolean {
    return this.token !== null;
  }

  /**
   * Get authorization headers if token is set
   */
  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Accept': 'application/vnd.github.v3+json',
    };

    if (this.token) {
      headers['Authorization'] = `token ${this.token}`;
    }

    return headers;
  }

  /**
   * Make a request to GitHub API
   */
  private async request<T>(endpoint: string): Promise<T> {
    const url = `${GITHUB_API_BASE}${endpoint}`;
    const response = await fetch(url, {
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      const error = await response.json() as GithubApiError;
      throw new Error(error.message || `GitHub API error: ${response.status}`);
    }

    return response.json() as Promise<T>;
  }

  /**
   * Get repository information
   */
  async getRepository(owner: string, repo: string): Promise<GithubRepository> {
    return this.request<GithubRepository>(`/repos/${owner}/${repo}`);
  }

  /**
   * Get branch information
   */
  async getBranch(owner: string, repo: string, branch: string): Promise<GithubBranch> {
    return this.request<GithubBranch>(`/repos/${owner}/${repo}/branches/${branch}`);
  }

  /**
   * Get latest commits for a repository
   */
  async getLatestCommits(owner: string, repo: string, branch: string, count = 1): Promise<GithubCommit[]> {
    return this.request<GithubCommit[]>(`/repos/${owner}/${repo}/commits?sha=${branch}&per_page=${count}`);
  }

  /**
   * Parse GitHub repository URL to extract owner and repo name
   * Handles formats like:
   * - https://github.com/owner/repo
   * - https://github.com/owner/repo.git
   * - git@github.com:owner/repo.git
   */
  parseRepoUrl(url: string): { owner: string; repo: string } | null {
    let match;
    
    // Format: https://github.com/owner/repo[.git]
    match = url.match(/github\.com[\/:]([\w-\.]+)\/([\w-\.]+)(\.git)?$/);
    
    if (match) {
      return {
        owner: match[1],
        repo: match[2]
      };
    }
    
    return null;
  }
}

// Export a singleton instance
export const githubApi = new GithubApiClient();