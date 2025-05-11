import { useSelector } from 'react-redux';
import { useGetAuthStatusQuery, useLoginMutation, useLogoutMutation } from '../api/authApi';
import { RootState } from '../store';

// This hook is a wrapper around the auth RTK Query endpoints
export const useAuth = () => {
  // Get current auth state from the store
  const { isAuthenticated, isAdmin, token } = useSelector((state: RootState) => state.auth);
  
  // RTK Query hooks
  const { refetch: checkAuthStatus } = useGetAuthStatusQuery(undefined, {
    // Skip initial fetch if we already have a token
    skip: false,
    refetchOnMountOrArgChange: true
  });
  
  const [loginMutation] = useLoginMutation();
  const [logoutMutation] = useLogoutMutation();

  const login = async (username: string, password: string) => {
    try {
      const result = await loginMutation({ username, password }).unwrap();
      return result.success;
    } catch (err) {
      console.error('Login failed:', err);
      return false;
    }
  };

  const logout = async () => {
    try {
      await logoutMutation().unwrap();
      return true;
    } catch (err) {
      console.error('Logout failed:', err);
      return false;
    }
  };

  return {
    isAuthenticated,
    isAdmin,
    token,
    login,
    logout,
    checkAuthStatus
  };
};