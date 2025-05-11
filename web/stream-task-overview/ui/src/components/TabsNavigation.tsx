import React from 'react';
import { LayoutGrid, FileText, PanelTop, ListTree, FileDigit } from 'lucide-react';
import { useAppSelector, useAppDispatch } from '../store/hooks';
import { changeTab } from '../store/slices/streamSlice';

const TabsNavigation: React.FC = () => {
  const { activeTab } = useAppSelector(state => state.stream);
  const dispatch = useAppDispatch();

  return (
    <div className="flex border-b-2 border-black mb-6 overflow-x-auto">
      <button
        onClick={() => dispatch(changeTab('main'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'main' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <LayoutGrid size={14} className="mr-2" />
        Dashboard
      </button>
      <button
        onClick={() => dispatch(changeTab('notes'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'notes' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <ListTree size={14} className="mr-2" />
        Notes
      </button>
      <button
        onClick={() => dispatch(changeTab('summary'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'summary' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <PanelTop size={14} className="mr-2" />
        Summary
      </button>
      <button
        onClick={() => dispatch(changeTab('raw'))}
        className={`px-4 py-3 flex items-center text-sm uppercase tracking-wider ${activeTab === 'raw' ? 'bg-black text-white' : 'bg-white text-black hover:bg-gray-100'}`}
      >
        <FileDigit size={14} className="mr-2" />
        Transcript
      </button>
    </div>
  );
};

export default TabsNavigation;