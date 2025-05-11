import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../store';
import { useAuth } from '../hooks/useAuth';
import { fetchGithubInfo, connectToGithub } from '../store/slices/streamSlice';
import { useGetGithubInfoQuery, useConnectGithubMutation } from '../api/githubApi';
import {
  selectGithubData,
  selectGithubLoading,
  selectGithubError
} from '../store/selectors';

const GithubInfoPanel: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [token, setToken] = useState('');
  const [showToken, setShowToken] = useState(false);
  const { isAdmin } = useAuth();
  
  // RTK Query approach
  const { data: githubData, isLoading: isLoadingQuery } = useGetGithubInfoQuery(undefined, {
    // Skip initial fetch since we'll use the thunk
    skip: true
  });
  const [connectToGithubRepo, { isLoading: isConnecting }] = useConnectGithubMutation();
  
  // Use memoized selectors to prevent unnecessary re-renders
  const github = useSelector(selectGithubData);
  const loading = useSelector(selectGithubLoading);
  const error = useSelector(selectGithubError);

  // Fetch GitHub info on component mount
  useEffect(() => {
    if (github.isConnected) {
      dispatch(fetchGithubInfo());
    }
  }, [dispatch, github.isConnected]);

  const handleConnect = async () => {
    if (token.trim()) {
      // Parse the repo URL to get owner and repo name
      const repoUrlMatch = github.repoUrl.match(/github\.com[\/:]([\w-\.]+)\/([\w-\.]+)(?:\.git)?$/);
      
      if (repoUrlMatch) {
        const repoOwner = repoUrlMatch[1];
        const repoName = repoUrlMatch[2];
        
        // Option 1: Using RTK Query directly
        try {
          await connectToGithubRepo({
            token: token.trim(),
            repoOwner,
            repoName
          }).unwrap();
          setToken('');
          
          // Fetch updated GitHub info
          dispatch(fetchGithubInfo());
        } catch (err) {
          console.error('Failed to connect to GitHub:', err);
        }
        
        // Option 2: Using Redux thunk
        // dispatch(connectToGithub({
        //   token: token.trim(),
        //   repoOwner,
        //   repoName
        // }))
        //   .unwrap()
        //   .then(() => {
        //     setToken('');
        //     dispatch(fetchGithubInfo());
        //   })
        //   .catch(err => console.error('Failed to connect to GitHub:', err));
      } else {
        console.error('Invalid GitHub repository URL');
      }
    }
  };

  if (loading || isLoadingQuery) {
    return <div className="p-4 bg-white rounded-lg shadow">Loading GitHub information...</div>;
  }

  if (error && github.isConnected) {
    return (
      <div className="p-4 bg-red-100 text-red-800 rounded-lg shadow">
        Error loading GitHub information: {error}
        <button 
          onClick={() => dispatch(fetchGithubInfo())} 
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
        >
          Retry
        </button>
      </div>
    );
  }

  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">GitHub Integration</h2>
      
      {github.isConnected ? (
        <div>
          <div className="bg-green-50 p-3 rounded border border-green-200 mb-4">
            <span className="text-green-700 font-medium">✓ Connected to GitHub</span>
          </div>
          
          <div className="space-y-2">
            <div>
              <span className="text-gray-500">Repository: </span>
              <a href={github.repoUrl} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                {github.repoOwner}/{github.repoName}
              </a>
            </div>
            
            <div>
              <span className="text-gray-500">Branch: </span>
              <span>{github.currentBranch}</span>
            </div>
            
            {github.latestCommit.hash && (
              <div className="mt-4">
                <div className="font-medium mb-1">Latest Commit:</div>
                <div className="bg-gray-50 p-3 rounded border border-gray-200">
                  <div>
                    <a 
                      href={github.latestCommit.url} 
                      target="_blank" 
                      rel="noopener noreferrer" 
                      className="text-blue-500 hover:underline"
                    >
                      {github.latestCommit.hash.substring(0, 7)}
                    </a>
                    <span className="ml-2">{github.latestCommit.message}</span>
                  </div>
                  <div className="text-sm text-gray-500">
                    {github.latestCommit.author} - {new Date(github.latestCommit.date).toLocaleString()}
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      ) : (
        <div>
          <div className="mb-4">
            <p className="text-gray-700 mb-2">
              {isAdmin ? 'Connect to your GitHub repository to track commits and sync your progress.' : 'GitHub integration is not connected.'}
            </p>
            
            <div className="mb-4">
              <div className="text-sm text-gray-500 mb-1">Repository URL</div>
              <div className="font-medium">{github.repoUrl}</div>
              {isAdmin && (
                <div className="text-xs text-gray-500 mt-1">
                  Update this in the Stream Information section.
                </div>
              )}
            </div>
            
            {isAdmin && (
              <>
                <div className="mb-4">
                  <label className="text-sm text-gray-500 mb-1 block" htmlFor="github-token">
                    GitHub Personal Access Token
                  </label>
                  <div className="flex">
                    <input
                      id="github-token"
                      type={showToken ? "text" : "password"}
                      value={token}
                      onChange={(e) => setToken(e.target.value)}
                      className="flex-grow p-2 border border-gray-300 rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
                      placeholder="Enter your GitHub token"
                    />
                    <button
                      onClick={() => setShowToken(!showToken)}
                      className="bg-gray-200 px-3 border-t border-r border-b border-gray-300 rounded-r"
                    >
                      {showToken ? 'Hide' : 'Show'}
                    </button>
                  </div>
                  <div className="text-xs text-gray-500 mt-1">
                    Token needs repo scope permissions.
                  </div>
                </div>
                
                <button
                  onClick={handleConnect}
                  disabled={!token.trim() || isConnecting}
                  className="w-full bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded disabled:bg-gray-300"
                >
                  {isConnecting ? 'Connecting...' : 'Connect to GitHub'}
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default GithubInfoPanel;