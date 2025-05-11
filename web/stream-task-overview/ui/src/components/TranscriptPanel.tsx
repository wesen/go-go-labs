import React, { useState } from 'react';
import {
  Clock,
  Edit3,
  GitCommit,
  PlayCircle,
  CheckCircle,
  ExternalLink
} from 'lucide-react';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { addTranscriptNote } from '../store/slices/streamSlice';

const TranscriptPanel: React.FC = () => {
  const { transcript, isLoggedIn } = useAppSelector(state => state.stream);
  const dispatch = useAppDispatch();
  
  const [newNote, setNewNote] = useState('');
  
  const handleAddNote = () => {
    if (newNote.trim()) {
      dispatch(addTranscriptNote(newNote.trim()));
      setNewNote('');
    }
  };
  
  // Format date for display
  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };
  
  // Get icon based on transcript entry type
  const getEntryIcon = (type: string) => {
    switch(type) {
      case 'task_started':
        return <PlayCircle size={16} className="text-blue-600" />;
      case 'task_completed':
        return <CheckCircle size={16} className="text-green-600" />;
      case 'commit':
        return <GitCommit size={16} className="text-purple-600" />;
      case 'note':
        return <Edit3 size={16} className="text-gray-600" />;
      default:
        return <Clock size={16} className="text-gray-600" />;
    }
  };
  
  // Sort transcript by timestamp (most recent first)
  const sortedTranscript = [...transcript].sort((a, b) => 
    new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  );
  
  return (
    <div className="p-4 border-2 border-black">
      <div className="mb-6 flex justify-between items-center">
        <h2 className="uppercase tracking-wider font-bold">Stream Transcript</h2>
        <div className="text-xs text-gray-600">{sortedTranscript.length} entries</div>
      </div>
      
      {/* Add note form (for logged in users) */}
      {isLoggedIn && (
        <div className="mb-6 flex">
          <input
            type="text"
            value={newNote}
            onChange={(e) => setNewNote(e.target.value)}
            placeholder="Add note to transcript..."
            className="flex-grow p-2 border-2 border-black rounded-none bg-white text-black"
          />
          <button
            onClick={handleAddNote}
            className="px-4 py-2 bg-black text-white rounded-none hover:bg-gray-800 transition-colors uppercase text-xs tracking-wider"
          >
            Add Note
          </button>
        </div>
      )}
      
      {/* Transcript entries */}
      <div className="space-y-4">
        {sortedTranscript.map(entry => (
          <div key={entry.id} className="border-l-4 pl-4 py-1" 
               style={{ 
                 borderLeftColor: 
                   entry.type === 'task_started' ? '#2563eb' : 
                   entry.type === 'task_completed' ? '#16a34a' : 
                   entry.type === 'commit' ? '#9333ea' : '#6b7280'
               }}>
            <div className="flex items-center mb-1 text-xs text-gray-500">
              <Clock size={12} className="mr-1" />
              <span>{formatTime(entry.timestamp)}</span>
              <div className="ml-3 flex items-center">
                {getEntryIcon(entry.type)}
                <span className="ml-1 uppercase tracking-wider">
                  {entry.type.replace('_', ' ')}
                </span>
              </div>
            </div>
            
            <div className="text-sm">{entry.content}</div>
            
            {/* Show commit info for commit entries */}
            {entry.type === 'commit' && entry.commitHash && entry.commitUrl && (
              <div className="mt-1">
                <a
                  href={entry.commitUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-xs font-mono text-blue-900 hover:underline flex items-center"
                >
                  {entry.commitHash}
                  <ExternalLink size={10} className="ml-1" />
                </a>
              </div>
            )}
          </div>
        ))}
      </div>
      
      {/* Empty state */}
      {sortedTranscript.length === 0 && (
        <div className="text-center py-12 text-gray-500">
          <div className="text-sm">No transcript entries yet</div>
        </div>
      )}
    </div>
  );
};

export default TranscriptPanel;