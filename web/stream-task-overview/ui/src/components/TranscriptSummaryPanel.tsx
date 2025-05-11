import React from 'react';
import { useAppSelector } from '../store/hooks';
import { GithubInfo } from '../store/slices/streamSlice';
import { ExternalLink, GitCommit } from 'lucide-react';

const TranscriptSummaryPanel: React.FC = () => {
  const { transcript, github } = useAppSelector(state => state.stream);
  
  // Group entries by type and sort by timestamp
  const groupedEntries = transcript.reduce((acc, entry) => {
    if (!acc[entry.type]) {
      acc[entry.type] = [];
    }
    acc[entry.type].push(entry);
    return acc;
  }, {} as Record<string, typeof transcript>);
  
  // Sort each group by timestamp (newest first)
  Object.keys(groupedEntries).forEach(type => {
    groupedEntries[type].sort((a, b) => 
      new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
    );
  });
  
  // Format date for display
  const formatDate = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleDateString([], { 
      month: 'short', 
      day: 'numeric',
      hour: '2-digit', 
      minute: '2-digit' 
    });
  };
  
  // Generate task progression summary
  const generateTaskSummary = () => {
    const taskStarts = groupedEntries['task_started'] || [];
    const taskCompletions = groupedEntries['task_completed'] || [];
    
    if (taskStarts.length === 0) return null;
    
    const completedTasks = taskCompletions.map(task => task.taskName);
    const activeTask = taskStarts.find(task => !completedTasks.includes(task.taskName));
    
    return (
      <div className="mb-8">
        <h3 className="text-lg font-bold mb-3">Task Progression</h3>
        <p className="mb-3">
          During this stream, we've worked on {taskStarts.length} distinct tasks. 
          {taskCompletions.length > 0 && ` We have completed ${taskCompletions.length} tasks, `}
          {activeTask && ` and are currently working on "${activeTask.taskName}".`}
        </p>
        
        <div className="border-l-4 border-blue-900 pl-4 py-2 bg-blue-50 mb-4">
          <div className="font-medium">Current Task:</div>
          <div>{activeTask?.taskName || "No active task"}</div>
        </div>
        
        {taskCompletions.length > 0 && (
          <div>
            <div className="font-medium mb-2">Completed Tasks:</div>
            <ul className="list-disc pl-5 space-y-1">
              {taskCompletions.slice(0, 5).map(task => (
                <li key={task.id}>
                  {task.taskName} <span className="text-xs text-gray-500">({formatDate(task.timestamp)})</span>
                </li>
              ))}
              {taskCompletions.length > 5 && (
                <li className="text-gray-500">...and {taskCompletions.length - 5} more</li>
              )}
            </ul>
          </div>
        )}
      </div>
    );
  };
  
  // Generate code changes summary
  const generateCodeSummary = () => {
    const commits = groupedEntries['commit'] || [];
    
    if (commits.length === 0) return null;
    
    return (
      <div className="mb-8">
        <h3 className="text-lg font-bold mb-3">Code Changes</h3>
        <p className="mb-3">
          The stream includes {commits.length} code {commits.length === 1 ? 'commit' : 'commits'} to the repository
          {github.repoOwner && github.repoName && (
            <> <a 
              href={`https://github.com/${github.repoOwner}/${github.repoName}`}
              target="_blank"
              rel="noopener noreferrer"
              className="text-blue-600 hover:underline"
            >
              {github.repoOwner}/{github.repoName}
            </a></>
          )}.
          {github.currentBranch && ` We are working on the "${github.currentBranch}" branch.`}
        </p>
        
        <div className="space-y-4">
          {commits.slice(0, 3).map(commit => (
            <div key={commit.id} className="border-l-4 border-purple-800 pl-4 py-2 bg-purple-50">
              <div className="flex items-start">
                <GitCommit size={16} className="mt-1 mr-2 text-purple-800" />
                <div>
                  <div className="font-medium">{commit.content}</div>
                  <div className="text-xs text-gray-600 mt-1">
                    {commit.commitHash && (
                      <a 
                        href={commit.commitUrl} 
                        target="_blank"
                        rel="noopener noreferrer"
                        className="font-mono text-blue-600 hover:underline flex items-center inline-block mr-2"
                      >
                        {commit.commitHash}
                        <ExternalLink size={10} className="ml-1 inline" />
                      </a>
                    )}
                    <span>{formatDate(commit.timestamp)}</span>
                  </div>
                </div>
              </div>
            </div>
          ))}
          {commits.length > 3 && (
            <div className="text-sm text-gray-500 italic">
              ...and {commits.length - 3} more commits
            </div>
          )}
        </div>
      </div>
    );
  };
  
  // Generate notes summary
  const generateNotesSummary = () => {
    const notes = groupedEntries['note'] || [];
    
    if (notes.length === 0) return null;
    
    return (
      <div className="mb-8">
        <h3 className="text-lg font-bold mb-3">Stream Notes</h3>
        <p className="mb-3">
          There {notes.length === 1 ? 'is' : 'are'} {notes.length} {notes.length === 1 ? 'note' : 'notes'} from this stream session.
        </p>
        
        <div className="border-l-4 border-gray-500 pl-4 py-2 bg-gray-50 space-y-3">
          {notes.slice(0, 5).map(note => (
            <div key={note.id} className="border-b border-gray-200 pb-2 last:border-0 last:pb-0">
              <div className="text-sm">{note.content}</div>
              <div className="text-xs text-gray-500 mt-1">{formatDate(note.timestamp)}</div>
            </div>
          ))}
          {notes.length > 5 && (
            <div className="text-sm text-gray-500 italic">
              ...and {notes.length - 5} more notes
            </div>
          )}
        </div>
      </div>
    );
  };
  
  // Generate a summary paragraph of the entire stream
  const generateStreamSummary = () => {
    const taskStarts = groupedEntries['task_started'] || [];
    const taskCompletions = groupedEntries['task_completed'] || [];
    const commits = groupedEntries['commit'] || [];
    const notes = groupedEntries['note'] || [];
    
    // Find earliest timestamp to determine stream start
    const allTimestamps = transcript.map(entry => new Date(entry.timestamp).getTime());
    const earliestTimestamp = Math.min(...allTimestamps);
    const latestTimestamp = Math.max(...allTimestamps);
    
    // Calculate stream duration in minutes
    const durationMinutes = Math.round((latestTimestamp - earliestTimestamp) / (1000 * 60));
    
    return (
      <div className="mb-8 bg-black text-white p-4">
        <h3 className="text-lg font-bold mb-3 uppercase tracking-wider">Stream Summary</h3>
        <p className="mb-2">
          This {durationMinutes > 0 ? `${durationMinutes}-minute ` : ""}stream focused on 
          {taskStarts.length > 0 ? ` ${taskStarts.length} development tasks` : " development"}.
          {taskCompletions.length > 0 && ` We completed ${taskCompletions.length} tasks`}
          {commits.length > 0 && ` and made ${commits.length} code ${commits.length === 1 ? 'commit' : 'commits'}`}.
          {notes.length > 0 && ` Throughout the stream, ${notes.length} ${notes.length === 1 ? 'note was' : 'notes were'} recorded.`}
        </p>
        <p>
          {github.currentBranch && github.repoOwner && github.repoName && (
            <>Working on the <span className="font-mono">{github.currentBranch}</span> branch of {github.repoOwner}/{github.repoName}.</>
          )}
        </p>
      </div>
    );
  };
  
  return (
    <div className="p-4 border-2 border-black">
      <div className="mb-6">
        <h2 className="text-xl font-bold uppercase tracking-wider mb-2">Stream Transcript Summary</h2>
        <div className="text-sm text-gray-600">A structured summary of the stream activities and key events</div>
      </div>
      
      {generateStreamSummary()}
      {generateTaskSummary()}
      {generateCodeSummary()}
      {generateNotesSummary()}
      
      {transcript.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <div className="text-lg font-medium">No transcript data available</div>
          <div className="text-sm">Stream events will appear here as they occur</div>
        </div>
      )}
    </div>
  );
};

export default TranscriptSummaryPanel;