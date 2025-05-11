import React from 'react';
import { useAppSelector } from '../store/hooks';

const RawTranscriptPanel: React.FC = () => {
  const { transcript } = useAppSelector(state => state.stream);
  
  // Sort transcript by timestamp (oldest first to create a chronological narrative)
  const sortedTranscript = [...transcript].sort((a, b) => 
    new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
  );
  
  // Format time function
  const formatTime = (timestamp: string) => {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  };
  
  // Format duration (seconds to MM:SS)
  const formatDuration = (seconds?: number) => {
    if (!seconds) return '';
    const minutes = Math.floor(seconds / 60);
    const remainingSeconds = seconds % 60;
    return `${minutes}:${remainingSeconds.toString().padStart(2, '0')}`;
  };
  
  // Filter the transcript entries
  const transcriptEntries = sortedTranscript.filter(entry => entry.type === 'transcript');
  
  // Check if we have transcript entries
  const hasTranscriptEntries = transcriptEntries.length > 0;
  
  return (
    <div className="p-4 border-2 border-black">
      <div className="mb-6">
        <h2 className="text-xl font-bold uppercase tracking-wider mb-2">Stream Transcript</h2>
        <div className="text-sm text-gray-600 mb-4">Raw transcript of the livestream conversation</div>
      </div>
      
      <div className="font-serif">
        {hasTranscriptEntries ? (
          <div className="space-y-4">
            {transcriptEntries.map(entry => (
              <div key={entry.id} className="pb-3 border-b border-gray-200 last:border-0">
                <div className="flex items-center justify-between text-sm text-gray-700 mb-1">
                  <div className="font-medium">{entry.speaker}</div>
                  <div className="flex items-center space-x-2">
                    <span>{formatTime(entry.timestamp)}</span>
                    {entry.duration && (
                      <span className="text-xs text-gray-500">({formatDuration(entry.duration)})</span>
                    )}
                  </div>
                </div>
                <div className="leading-relaxed">{entry.content}</div>
              </div>
            ))}
          </div>
        ) : (
          <div className="text-center py-12 text-gray-500">
            <div className="text-lg font-medium">No transcript available</div>
            <div className="text-sm">Transcript entries will appear here once they are recorded</div>
          </div>
        )}
      </div>
    </div>
  );
};

export default RawTranscriptPanel;