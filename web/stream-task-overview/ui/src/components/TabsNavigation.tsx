import React from 'react';
import { LayoutGrid, FileText } from 'lucide-react';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { changeTab } from '../store/slices/streamSlice';

const TabsNavigation: React.FC = () => {
  const { activeTab } = useAppSelector(state => state.stream);
  const dispatch = useAppDispatch();

  return (
    <div className="flex border-b-2 border-black mb-6">
      <button
        onClick={() => dispatch(changeTab('main'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'main' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <LayoutGrid size={14} className="mr-2" />
        Main Dashboard
      </button>
      <button
        onClick={() => dispatch(changeTab('transcript'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'transcript' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <FileText size={14} className="mr-2" />
        Stream Transcript
      </button>
    </div>
  );
};

export default TabsNavigation;