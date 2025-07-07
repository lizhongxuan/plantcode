import React from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { authApi } from '@/services/api';
import Button from '@/components/ui/Button';
import { User, LogOut, Settings } from 'lucide-react';

const Header: React.FC = () => {
  const navigate = useNavigate();
  const { user, clearAuth } = useAuthStore();

  const handleLogout = () => {
    authApi.logout();
    clearAuth();
    navigate('/login');
  };

  return (
    <header className="main-header">
      <div className="h-full px-4 flex items-center justify-between">
        <div className="flex items-center space-x-4">
          <h1 className="text-xl font-semibold text-gray-900">
            AI开发平台
          </h1>
        </div>
        
        <div className="flex items-center space-x-4">
          {user && (
            <div className="flex items-center space-x-3">
              <div className="flex items-center space-x-2">
                <User className="h-5 w-5 text-gray-500" />
                <span className="text-sm text-gray-700">{user.full_name}</span>
              </div>
              
              <div className="flex items-center space-x-2">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => navigate('/settings')}
                >
                  <Settings className="h-4 w-4" />
                </Button>
                
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={handleLogout}
                >
                  <LogOut className="h-4 w-4" />
                </Button>
              </div>
            </div>
          )}
        </div>
      </div>
    </header>
  );
};

export default Header; 