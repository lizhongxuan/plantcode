import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import PrivateRoute from '../PrivateRoute';
import { useAuthStore } from '@/store';
import { User } from '@/types';

// Mock the auth store
jest.mock('@/store', () => ({
  useAuthStore: jest.fn(),
}));

// Mock Navigate component
jest.mock('react-router-dom', () => ({
  ...jest.requireActual('react-router-dom'),
  Navigate: ({ to }: { to: string }) => <div data-testid="navigate">{to}</div>,
}));

const mockUseAuthStore = useAuthStore as jest.MockedFunction<typeof useAuthStore>;

const TestComponent = () => <div data-testid="protected-content">Protected Content</div>;

describe('PrivateRoute', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('should render children when user is authenticated', () => {
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

    mockUseAuthStore.mockReturnValue({
      isAuthenticated: true,
      user: mockUser,
      token: 'jwt.token.here',
      setAuth: jest.fn(),
      clearAuth: jest.fn(),
      updateUser: jest.fn(),
    });

    // Act
    render(
      <BrowserRouter>
        <PrivateRoute>
          <TestComponent />
        </PrivateRoute>
      </BrowserRouter>
    );

    // Assert
    expect(screen.getByTestId('protected-content')).toBeInTheDocument();
    expect(screen.queryByTestId('navigate')).not.toBeInTheDocument();
  });

  it('should redirect to login when user is not authenticated', () => {
    // Arrange
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false,
      user: null,
      token: null,
      setAuth: jest.fn(),
      clearAuth: jest.fn(),
      updateUser: jest.fn(),
    });

    // Act
    render(
      <BrowserRouter>
        <PrivateRoute>
          <TestComponent />
        </PrivateRoute>
      </BrowserRouter>
    );

    // Assert
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
    expect(screen.getByTestId('navigate')).toHaveTextContent('/login');
  });

  it('should redirect to login when user is null even if isAuthenticated is true', () => {
    // Arrange
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: true,
      user: null, // 用户为空
      token: 'jwt.token.here',
      setAuth: jest.fn(),
      clearAuth: jest.fn(),
      updateUser: jest.fn(),
    });

    // Act
    render(
      <BrowserRouter>
        <PrivateRoute>
          <TestComponent />
        </PrivateRoute>
      </BrowserRouter>
    );

    // Assert
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
    expect(screen.getByTestId('navigate')).toHaveTextContent('/login');
  });

  it('should redirect to login when isAuthenticated is false even if user exists', () => {
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

    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false, // 未认证
      user: mockUser,
      token: null,
      setAuth: jest.fn(),
      clearAuth: jest.fn(),
      updateUser: jest.fn(),
    });

    // Act
    render(
      <BrowserRouter>
        <PrivateRoute>
          <TestComponent />
        </PrivateRoute>
      </BrowserRouter>
    );

    // Assert
    expect(screen.queryByTestId('protected-content')).not.toBeInTheDocument();
    expect(screen.getByTestId('navigate')).toHaveTextContent('/login');
  });
}); 