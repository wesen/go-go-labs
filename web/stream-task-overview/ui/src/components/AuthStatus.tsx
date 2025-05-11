import React, { useState, useEffect } from 'react';
import { User, LogOut, Shield } from 'lucide-react';

interface AuthUser {
  authenticated: boolean;
  id?: string;
  login?: string;
  name?: string;
  email?: string;
  avatar?: string;
  is_admin?: boolean;
}

const AuthStatus: React.FC = () => {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await fetch('/auth/user');
        const data = await response.json();
        setUser(data);
      } catch (error) {
        console.error('Error fetching user data:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, []);

  if (loading) {
    return (
      <div className="p-2 border-2 border-black bg-gray-100">
        <div className="text-xs uppercase tracking-wider mb-1">AUTHENTICATION</div>
        <p className="text-sm">Loading...</p>
      </div>
    );
  }

  return (
    <div className="p-2 border-2 border-black bg-gray-100">
      <div className="text-xs uppercase tracking-wider mb-1">AUTHENTICATION</div>
      
      {user?.authenticated ? (
        <div className="flex flex-col">
          <div className="flex items-center">
            {user.avatar && (
              <img 
                src={user.avatar} 
                alt={user.name || user.login} 
                className="w-6 h-6 mr-2 border border-black"
              />
            )}
            <div className="flex-grow">
              <p className="font-bold">{user.name || user.login}</p>
              {user.login && user.name && <p className="text-xs">{user.login}</p>}
            </div>
            {user.is_admin && (
              <div className="w-6 h-6 mr-2 flex items-center justify-center bg-black text-white" title="Admin">
                <Shield size={14} />
              </div>
            )}
          </div>
          
          <div className="mt-2 text-right">
            <a 
              href="/auth/logout" 
              className="px-3 py-1 bg-black text-white text-xs uppercase tracking-wider inline-flex items-center"
            >
              <LogOut size={12} className="mr-1" /> Logout
            </a>
          </div>
        </div>
      ) : (
        <div className="flex flex-col items-center">
          <p className="text-sm mb-2">Not logged in</p>
          <a 
            href="/auth/github" 
            className="px-3 py-1 bg-black text-white text-xs uppercase tracking-wider inline-flex items-center"
          >
            <User size={12} className="mr-1" /> Login with GitHub
          </a>
        </div>
      )}
    </div>
  );
};

export default AuthStatus;