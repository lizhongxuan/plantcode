import React from 'react';
import { NavLink } from 'react-router-dom';
import { clsx } from 'clsx';
import { 
  Home, 
  FolderOpen, 
  Plus, 
  Settings
} from 'lucide-react';

interface SidebarItem {
  to: string;
  icon: React.ReactNode;
  label: string;
}

const sidebarItems: SidebarItem[] = [
  {
    to: '/dashboard',
    icon: <Home className="h-5 w-5" />,
    label: '仪表板',
  },
  {
    to: '/projects',
    icon: <FolderOpen className="h-5 w-5" />,
    label: '项目管理',
  },
  {
    to: '/projects/new',
    icon: <Plus className="h-5 w-5" />,
    label: '新建项目',
  },
  {
    to: '/settings',
    icon: <Settings className="h-5 w-5" />,
    label: '设置',
  },
];

const Sidebar: React.FC = () => {
  return (
    <div className="main-sidebar">
      <nav className="flex-1 py-4">
        <div className="px-4 mb-4">
          <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider">
            导航
          </h2>
        </div>
        
        <ul className="space-y-1 px-2">
          {sidebarItems.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                className={({ isActive }) =>
                  clsx(
                    'sidebar-nav-item',
                    isActive ? 'sidebar-nav-item-active' : 'sidebar-nav-item-inactive'
                  )
                }
              >
                {item.icon}
                <span className="ml-3">{item.label}</span>
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </div>
  );
};

export default Sidebar; 