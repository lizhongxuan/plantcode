import React from 'react';
import { NavLink } from 'react-router-dom';
import { clsx } from 'clsx';
import { 
  Home, 
  FolderOpen, 
  Plus, 
  Settings, 
  FileText,
  MessageSquare 
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
    to: '/documents',
    icon: <FileText className="h-5 w-5" />,
    label: '文档管理',
  },
  {
    to: '/chat',
    icon: <MessageSquare className="h-5 w-5" />,
    label: 'AI对话',
  },
  {
    to: '/settings',
    icon: <Settings className="h-5 w-5" />,
    label: '设置',
  },
];

const Sidebar: React.FC = () => {
  return (
    <div className="bg-white border-r border-gray-200 w-64 min-h-screen">
      <nav className="mt-4">
        <div className="px-4">
          <h2 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-3">
            导航
          </h2>
        </div>
        
        <ul className="space-y-1">
          {sidebarItems.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                className={({ isActive }) =>
                  clsx(
                    'flex items-center px-4 py-2 text-sm font-medium transition-colors duration-200',
                    isActive
                      ? 'bg-primary-50 text-primary-700 border-r-2 border-primary-600'
                      : 'text-gray-700 hover:bg-gray-50 hover:text-gray-900'
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