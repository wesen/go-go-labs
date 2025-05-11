import React, { useState, useEffect } from 'react';
import { Github, GitBranch, GitCommit, AlertCircle, ExternalLink } from 'lucide-react';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { 
  setGithubToken,
  setGithubConnectionStatus,
  setGithubError,
  updateGithubBranchInfo,
  updateGithubCommitInfo
} from '../store/slices/streamSlice';
import { githubApi } from '../services/githubApi';

const GithubInfoPanel: React.FC = () => {
  const { github, isLoggedIn } = useAppSelector(state => state.stream);
  const dispatch = useAppDispatch();
  
  const [tokenInput, setTokenInput] = useState('');
  const [isTokenModalOpen, setIsTokenModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  
  const handleTokenSubmit = async () => {
    if (!tokenInput.trim()) return;
    
    setIsLoading(true);
    try {
      // Set token in the API client
      githubApi.setToken(tokenInput.trim());
      
      // Try to fetch repo info to validate token
      if (github.repoOwner && github.repoName) {
        await githubApi.getRepository(github.repoOwner, github.repoName);
        
        // If we get here, token is valid
        dispatch(setGithubToken(tokenInput.trim()));
        dispatch(setGithubConnectionStatus(true));
        dispatch(setGithubError(undefined));
        setIsTokenModalOpen(false);
        
        // Fetch initial data
        fetchGithubData();
      }
    } catch (error) {
      dispatch(setGithubError(error.message || 'Failed to connect to GitHub'));
      dispatch(setGithubConnectionStatus(false));
    } finally {
      setIsLoading(false);
    }
  };
  
  const fetchGithubData = async () => {
    if (!github.isConnected || !github.repoOwner || !github.repoName) return;
    
    setIsLoading(true);
    try {
      // Fetch default branch if not specified
      let branch = github.currentBranch;
      if (!branch) {
        const repo = await githubApi.getRepository(github.repoOwner, github.repoName);
        branch = repo.default_branch;
        dispatch(updateGithubBranchInfo({ branch }));
      }
      
      // Fetch latest commit
      const commits = await githubApi.getLatestCommits(github.repoOwner, github.repoName, branch, 1);
      if (commits.length > 0) {
        const latestCommit = commits[0];
        dispatch(updateGithubCommitInfo({
          message: latestCommit.commit.message,
          author: latestCommit.author?.login || latestCommit.commit.author.name,
          hash: latestCommit.sha.substring(0, 7),
          date: latestCommit.commit.author.date,
          url: latestCommit.html_url
        }));
      }
    } catch (error) {
      dispatch(setGithubError(error.message || 'Failed to fetch GitHub data'));
    } finally {
      setIsLoading(false);
    }
  };
  
  // Initial data fetch and polling setup
  useEffect(() => {
    if (github.isConnected && github.token) {
      // Set token in API client
      githubApi.setToken(github.token);
      
      // Initial fetch
      fetchGithubData();
      
      // Poll every minute
      const interval = setInterval(fetchGithubData, 60000);
      return () => clearInterval(interval);
    }
  }, [github.isConnected, github.token, github.repoOwner, github.repoName]);
  
  // Auto show token modal when repo info exists but no token
  useEffect(() => {
    if (github.repoOwner && github.repoName && !github.token && !isTokenModalOpen && !github.isConnected) {
      setIsTokenModalOpen(true);
    }
  }, [github.repoOwner, github.repoName, github.token, github.isConnected]);
  
  // Disconnect from GitHub
  const handleDisconnect = () => {
    githubApi.clearToken();
    dispatch(setGithubToken(''));
    dispatch(setGithubConnectionStatus(false));
  };
  
  // Format commit date
  const formatDate = (dateString: string) => {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleString();
  };
  
  return (
    <div className="border-2 border-black p-4 mt-4">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center">
          <div className="w-6 h-6 mr-2 flex items-center justify-center bg-black text-white">
            <Github size={16} />
          </div>
          <div className="text-xs uppercase tracking-wider">GITHUB REPOSITORY</div>
        </div>
        
        {isLoggedIn && (
          github.isConnected ? (
            <button
              onClick={handleDisconnect}
              className="px-3 py-1 bg-red-900 text-white rounded-none hover:bg-red-800 transition-colors uppercase text-xs tracking-wider"
            >
              Disconnect
            </button>
          ) : (
            <button
              onClick={() => setIsTokenModalOpen(true)}
              className="px-3 py-1 bg-black text-white rounded-none hover:bg-gray-800 transition-colors uppercase text-xs tracking-wider"
            >
              Connect
            </button>
          )
        )}
      </div>
      
      {/* Repository info */}
      <div className="mb-4">
        <a 
          href={github.repoUrl} 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-blue-900 hover:underline flex items-center"
        >
          {github.repoOwner}/{github.repoName}
          <ExternalLink size={14} className="ml-1" />
        </a>
      </div>
      
      {/* Loading or error state */}
      {isLoading && <div className="text-xs italic mb-3">Loading GitHub data...</div>}
      {github.error && (
        <div className="text-red-800 text-xs mb-3 flex items-center">
          <AlertCircle size={12} className="mr-1" /> {github.error}
        </div>
      )}
      
      {/* Connected state with branch and commit info */}
      {github.isConnected && (
        <div className="space-y-3">
          <div className="flex items-center">
            <div className="w-5 h-5 mr-1 flex items-center justify-center">
              <GitBranch size={14} />
            </div>
            <div>
              <div className="text-xs uppercase tracking-wider">BRANCH</div>
              <span className="text-sm">{github.currentBranch}</span>
            </div>
          </div>
          
          {github.latestCommit.hash && (
            <div className="border-t border-gray-200 pt-3">
              <div className="flex items-start">
                <div className="w-5 h-5 mr-1 flex items-center justify-center mt-0.5">
                  <GitCommit size={14} />
                </div>
                <div>
                  <div className="text-xs uppercase tracking-wider">LATEST COMMIT</div>
                  <div className="text-sm font-medium">{github.latestCommit.message}</div>
                  <div className="text-xs mt-1">
                    <a 
                      href={github.latestCommit.url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="font-mono text-blue-900 hover:underline"
                    >
                      {github.latestCommit.hash}
                    </a>
                    {' '} by {github.latestCommit.author} at {formatDate(github.latestCommit.date)}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      )}
      
      {/* Token input modal */}
      {isTokenModalOpen && (
        <div className="border-2 border-black p-4 mt-4 bg-gray-50">
          <div className="text-sm font-medium mb-2">Enter GitHub Personal Access Token</div>
          <div className="text-xs mb-3">
            Create a token with <code>repo</code> scope at{' '}
            <a 
              href="https://github.com/settings/tokens/new" 
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-900 hover:underline"
            >
              github.com/settings/tokens/new
            </a>
          </div>
          <input
            type="password"
            value={tokenInput}
            onChange={(e) => setTokenInput(e.target.value)}
            placeholder="ghp_xxxxxxxxxxxx"
            className="w-full p-2 border-2 border-black rounded-none bg-white text-black mb-2"
          />
          <div className="flex space-x-2">
            <button
              onClick={handleTokenSubmit}
              disabled={isLoading}
              className="px-3 py-1 bg-black text-white rounded-none hover:bg-gray-800 transition-colors uppercase text-xs tracking-wider"
            >
              {isLoading ? 'Connecting...' : 'Connect'}
            </button>
            <button
              onClick={() => setIsTokenModalOpen(false)}
              className="px-3 py-1 bg-gray-500 text-white rounded-none hover:bg-gray-400 transition-colors uppercase text-xs tracking-wider"
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default GithubInfoPanel;