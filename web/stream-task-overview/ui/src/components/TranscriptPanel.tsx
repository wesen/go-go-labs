import React, { useEffect, useState, useMemo } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { AppDispatch } from '../store';
import { TranscriptEntry, fetchTranscript, addNote } from '../store/slices/streamSlice';
import { useGetTranscriptQuery, useAddNoteMutation } from '../api/transcriptApi';
import { 
  selectTranscriptData, 
  selectTranscriptLoading, 
  selectTranscriptError 
} from '../store/selectors';

const TranscriptPanel: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const [newNote, setNewNote] = useState('');
  
  // Use memoized selectors to prevent unnecessary re-renders
  const transcript = useSelector(selectTranscriptData);
  const loading = useSelector(selectTranscriptLoading);
  const error = useSelector(selectTranscriptError);
  
  // RTK Query approach
  const { data: transcriptData, isLoading: isLoadingQuery } = useGetTranscriptQuery(undefined, {
    // Skip initial fetch since we'll use the thunk
    skip: true
  });
  const [addNoteToTranscript, { isLoading: isAddingNote }] = useAddNoteMutation();

  // Fetch transcript on component mount
  useEffect(() => {
    dispatch(fetchTranscript());
  }, [dispatch]);

  const handleAddNote = async () => {
    if (newNote.trim()) {
      // Option 1: Using RTK Query directly
      try {
        await addNoteToTranscript(newNote.trim()).unwrap();
        setNewNote('');
      } catch (err) {
        console.error('Failed to add note:', err);
      }
    }
  };

  const formatTime = (timestamp: string) => {
    try {
      const date = new Date(timestamp);
      return date.toLocaleTimeString();
    } catch (e) {
      return 'Invalid time';
    }
  };

  // Use a memoized version of the sorted transcript - moved up before conditionals
  const sortedTranscript = useMemo(() => {
    if (!transcript || transcript.length === 0) return [];
    return [...transcript].sort((a, b) => {
      const dateA = new Date(a.timestamp).getTime();
      const dateB = new Date(b.timestamp).getTime();
      return dateB - dateA; // Sort in descending order (newest first)
    });
  }, [transcript]);
  
  const renderContent = () => {
    if (loading || isLoadingQuery) {
      return <div className="p-4 bg-white rounded-lg shadow">Loading transcript...</div>;
    }

    if (error) {
      return (
        <div className="p-4 bg-red-100 text-red-800 rounded-lg shadow">
          Error loading transcript: {error}
          <button 
            onClick={() => dispatch(fetchTranscript())} 
            className="mt-2 px-4 py-2 bg-blue-500 text-white rounded"
          >
            Retry
          </button>
        </div>
      );
    }
    
    return (
      <div className="space-y-4">
        {!sortedTranscript || sortedTranscript.length === 0 ? (
          <p className="text-gray-500 italic">No transcript entries yet</p>
        ) : (
          sortedTranscript.map((entry) => (
            <div key={entry.id} className="flex">
              <div className="w-20 flex-shrink-0 text-sm text-gray-500">
                {formatTime(entry.timestamp)}
              </div>
              <div className="flex-grow">
                {renderEntryContent(entry)}
              </div>
            </div>
          ))
        )}
      </div>
    );
  }

  const renderEntryContent = (entry: TranscriptEntry) => {
    switch (entry.type) {
      case 'task_started':
        return (
          <div className="bg-blue-50 p-3 rounded border border-blue-200">
            <div className="font-medium text-blue-700">▶ Started: {entry.taskName}</div>
            <div>{entry.content}</div>
          </div>
        );
        
      case 'task_completed':
        return (
          <div className="bg-green-50 p-3 rounded border border-green-200">
            <div className="font-medium text-green-700">✓ Completed: {entry.taskName}</div>
            <div>{entry.content}</div>
          </div>
        );
        
      case 'commit':
        return (
          <div className="bg-purple-50 p-3 rounded border border-purple-200">
            <div className="font-medium">
              <span className="text-purple-700">📝 Commit:</span> 
              {entry.commitUrl ? (
                <a href={entry.commitUrl} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline ml-1">
                  {entry.commitHash?.substring(0, 7)}
                </a>
              ) : (
                <span className="ml-1">{entry.commitHash?.substring(0, 7)}</span>
              )}
            </div>
            <div>{entry.content}</div>
          </div>
        );
        
      case 'note':
        return (
          <div className="bg-yellow-50 p-3 rounded border border-yellow-200">
            <div className="font-medium text-yellow-700">📝 Note</div>
            <div>{entry.content}</div>
          </div>
        );
        
      case 'paragraph':
        return (
          <div className="bg-gray-50 p-3 rounded border border-gray-200">
            {entry.title && <div className="font-medium text-gray-700 mb-2">{entry.title}</div>}
            <div className="whitespace-pre-line">{entry.content}</div>
            {entry.timeRange && (
              <div className="text-xs text-gray-500 mt-2">
                {formatTime(entry.timeRange.start)} - {formatTime(entry.timeRange.end)}
              </div>
            )}
          </div>
        );
        
      case 'transcript':
        return (
          <div className="p-3 rounded border border-gray-200">
            <div className="font-medium">{entry.speaker}</div>
            <div>"{entry.content}"</div>
            {entry.duration && (
              <div className="text-xs text-gray-500">{entry.duration}s</div>
            )}
          </div>
        );
        
      default:
        return <div>{entry.content}</div>;
    }
  };

  return (
    <div className="p-4 bg-white rounded-lg shadow">
      <h2 className="text-xl font-bold mb-4">Transcript & Notes</h2>
      
      {/* Add Note Form */}
      <div className="mb-6">
        <div className="flex">
          <input
            type="text"
            value={newNote}
            onChange={(e) => setNewNote(e.target.value)}
            className="flex-grow p-2 border border-gray-300 rounded-l focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Add a note..."
          />
          <button
            onClick={handleAddNote}
            disabled={!newNote.trim() || isAddingNote}
            className="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-r disabled:bg-gray-300"
          >
            {isAddingNote ? 'Adding...' : 'Add Note'}
          </button>
        </div>
      </div>
      
      {/* Transcript Entries */}
      {renderContent()}
    </div>
  );
};

export default TranscriptPanel;