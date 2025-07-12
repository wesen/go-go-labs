import React from 'react';
import { WidgetList } from './WidgetList.tsx';
import { FeaturedWidget } from './FeaturedWidget.tsx';

export const WidgetManager: React.FC = () => {
  return (
    <div className="widget-manager">
      <div className="widget-manager-left">
        <h2>Available Widgets</h2>
        <WidgetList />
      </div>
      <div className="widget-manager-right">
        <FeaturedWidget />
      </div>
    </div>
  );
}; 