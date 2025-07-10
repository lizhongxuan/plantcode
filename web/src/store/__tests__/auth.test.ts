import { useAuthStore } from '../index';
import { User } from '@/types';

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('Auth Store', () => {
  beforeEach(() => {
    // Reset store
    useAuthStore.setState({
      isAuthenticated: false,
      user: null,
      token: null,
    });
    jest.clearAllMocks();
  });

  describe('setAuth', () => {
    it('should set authentication state correctly', () => {
      // Arrange
      const mockUser: User = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login: new Date().toISOString(),
        preferences: '{}',
      };
      const mockToken = 'jwt.token.here';

      // Act
      useAuthStore.getState().setAuth(mockUser, mockToken);

      // Assert
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(true);
      expect(state.user).toEqual(mockUser);
      expect(state.token).toBe(mockToken);
      expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', mockToken);
    });
  });

  describe('clearAuth', () => {
    it('should clear authentication state correctly', () => {
      // Arrange - 先设置认证状态
      const mockUser: User = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login: new Date().toISOString(),
        preferences: '{}',
      };
      const mockToken = 'jwt.token.here';
      
      useAuthStore.getState().setAuth(mockUser, mockToken);

      // Act
      useAuthStore.getState().clearAuth();

      // Assert
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
      expect(state.token).toBeNull();
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token');
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth-store');
    });
  });

  describe('updateUser', () => {
    it('should update user information correctly', () => {
      // Arrange
      const initialUser: User = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login: new Date().toISOString(),
        preferences: '{}',
      };
      
      const updatedUser: User = {
        ...initialUser,
        full_name: 'Updated Test User',
        email: 'updated@example.com',
      };

      useAuthStore.getState().setAuth(initialUser, 'jwt.token.here');

      // Act
      useAuthStore.getState().updateUser(updatedUser);

      // Assert
      const state = useAuthStore.getState();
      expect(state.user).toEqual(updatedUser);
      expect(state.isAuthenticated).toBe(true); // 应该保持认证状态
      expect(state.token).toBe('jwt.token.here'); // token应该保持不变
    });

    it('should not affect authentication state when updating user', () => {
      // Arrange
      const initialUser: User = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
        created_at: new Date().toISOString(),
        updated_at: new Date().toISOString(),
        last_login: new Date().toISOString(),
        preferences: '{}',
      };

      useAuthStore.getState().setAuth(initialUser, 'jwt.token.here');
      const initialState = useAuthStore.getState();

      // Act
      useAuthStore.getState().updateUser({
        ...initialUser,
        full_name: 'Updated Name',
      });

      // Assert
      const finalState = useAuthStore.getState();
      expect(finalState.isAuthenticated).toBe(initialState.isAuthenticated);
      expect(finalState.token).toBe(initialState.token);
    });
  });

  describe('initial state', () => {
    it('should have correct initial state', () => {
      // Act
      const state = useAuthStore.getState();

      // Assert
      expect(state.isAuthenticated).toBe(false);
      expect(state.user).toBeNull();
      expect(state.token).toBeNull();
    });
  });
}); 