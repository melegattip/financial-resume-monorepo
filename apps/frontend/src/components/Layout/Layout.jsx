import React, { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';
import Header from './Header';
import { PeriodProvider } from '../../contexts/PeriodContext';

const Layout = () => {
  const [isDesktopCollapsed, setIsDesktopCollapsed] = useState(false);
  const [isMobileOpen, setIsMobileOpen] = useState(false);

  return (
    <PeriodProvider>
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
        {/* Sidebar */}
        <Sidebar
          isDesktopCollapsed={isDesktopCollapsed}
          onDesktopToggle={() => setIsDesktopCollapsed(prev => !prev)}
          isMobileOpen={isMobileOpen}
          onMobileClose={() => setIsMobileOpen(false)}
        />

        {/* Main content area */}
        <div className={`transition-all duration-300 ease-in-out ${isDesktopCollapsed ? 'lg:ml-16' : 'lg:ml-52'}`}>
          {/* Header */}
          <Header onMobileMenuToggle={() => setIsMobileOpen(prev => !prev)} isMobileMenuOpen={isMobileOpen} />

          {/* Page content */}
          <main className="p-3 sm:p-4 pt-3">
            <Outlet />
          </main>
        </div>
      </div>
    </PeriodProvider>
  );
};

export default Layout;
