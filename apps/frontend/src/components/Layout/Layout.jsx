import React from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import { PeriodProvider } from '../../contexts/PeriodContext';

const Layout = () => {
  return (
    <PeriodProvider>
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
        {/* Sidebar */}
        <Sidebar />
        
        {/* Main content area */}
        <div className="lg:ml-64">
          {/* Header */}
          <Header />
          
          {/* Page content */}
          <main className="p-6 pt-4 sm:pt-6">
            <Outlet />
          </main>
        </div>
      </div>
    </PeriodProvider>
  );
};

export default Layout; 