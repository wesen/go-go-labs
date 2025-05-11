import React from 'react';
import { User, LogOut, Shield } from 'lucide-react';
import { useAuth } from '../hooks/useAuth';

const AuthStatus: React.FC = () => {
  const { isAuthenticated, isAdmin, logout } = useAuth();
  
  const handleLogout = async () => {
    await logout();
    // Redirect to login page if needed
    // window.location.href = '/';
  };

  // Loading state is managed by RTK Query

  return (
    <div className="p-2 border-2 border-black bg-gray-100">
      <div className="text-xs uppercase tracking-wider mb-1">AUTHENTICATION</div>
      
      {isAuthenticated ? (
        <div className="flex flex-col">
          <div className="flex items-center">
            <div className="flex-grow">
              <p className="font-bold">Authenticated User</p>
            </div>
            {isAdmin && (
              <div className="w-6 h-6 mr-2 flex items-center justify-center bg-black text-white" title="Admin">
                <Shield size={14} />
              </div>
            )}
          </div>
          
          <div className="mt-2 text-right">
            <button 
              onClick={handleLogout}
              className="px-3 py-1 bg-black text-white text-xs uppercase tracking-wider inline-flex items-center"
            >
              <LogOut size={12} className="mr-1" /> Logout
            </button>
          </div>
        </div>
      ) : (
        <div className="flex flex-col items-center">
          <p className="text-sm mb-2">Not logged in</p>
          <a 
            href="/auth/login" 
            className="px-3 py-1 bg-black text-white text-xs uppercase tracking-wider inline-flex items-center"
          >
            <User size={12} className="mr-1" /> Login
          </a>
        </div>
      )}
    </div>
  );
};

export default AuthStatus;