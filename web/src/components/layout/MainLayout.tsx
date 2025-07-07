import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import Sidebar from './Sidebar';

const MainLayout: React.FC = () => {
  return (
    <div className="main-layout">
      <Header />
      <div className="flex flex-1 h-[calc(100vh-4rem)]">
        <aside className="flex-shrink-0">
          <Sidebar />
        </aside>
        <main className="main-content">
          <div className="p-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};

export default MainLayout; 