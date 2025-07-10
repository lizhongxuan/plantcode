import { authApi } from '../api';
import axios from 'axios';
import { LoginRequest, CreateUserRequest } from '@/types';

// Mock axios
jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

// Mock localStorage
const localStorageMock = {
  getItem: jest.fn(),
  setItem: jest.fn(),
  removeItem: jest.fn(),
  clear: jest.fn(),
};
Object.defineProperty(window, 'localStorage', { value: localStorageMock });

describe('AuthAPI', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    localStorageMock.getItem.mockClear();
    localStorageMock.setItem.mockClear();
    localStorageMock.removeItem.mockClear();
  });

  describe('register', () => {
    it('should register user successfully', async () => {
      // Arrange
      const userData: CreateUserRequest = {
        username: 'testuser',
        email: 'test@example.com',
        password: 'password123',
        full_name: 'Test User',
      };

      const expectedUser = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
      };

      const mockResponse = {
        data: {
          success: true,
          data: expectedUser,
          message: '用户注册成功',
        },
      };

      mockedAxios.post.mockResolvedValue(mockResponse);

      // Act
      const result = await authApi.register(userData);

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith('/auth/register', userData);
      expect(result).toEqual(expectedUser);
    });

    it('should throw error when registration fails', async () => {
      // Arrange
      const userData: CreateUserRequest = {
        username: 'testuser',
        email: 'test@example.com',
        password: 'password123',
        full_name: 'Test User',
      };

      const mockResponse = {
        data: {
          success: false,
          error: '邮箱已被注册',
        },
      };

      mockedAxios.post.mockResolvedValue(mockResponse);

      // Act & Assert
      await expect(authApi.register(userData)).rejects.toThrow('邮箱已被注册');
    });
  });

  describe('login', () => {
    it('should login user successfully', async () => {
      // Arrange
      const loginData: LoginRequest = {
        email: 'test@example.com',
        password: 'password123',
      };

      const expectedResponse = {
        user: {
          user_id: '123',
          username: 'testuser',
          email: 'test@example.com',
          full_name: 'Test User',
          status: 'active',
        },
        token: 'jwt.token.here',
      };

      const mockResponse = {
        data: {
          success: true,
          data: expectedResponse,
          message: '登录成功',
        },
      };

      mockedAxios.post.mockResolvedValue(mockResponse);

      // Act
      const result = await authApi.login(loginData);

      // Assert
      expect(mockedAxios.post).toHaveBeenCalledWith('/auth/login', loginData);
      expect(result).toEqual(expectedResponse);
      expect(localStorageMock.setItem).toHaveBeenCalledWith('auth_token', 'jwt.token.here');
    });

    it('should throw error when login fails', async () => {
      // Arrange
      const loginData: LoginRequest = {
        email: 'test@example.com',
        password: 'wrongpassword',
      };

      const mockResponse = {
        data: {
          success: false,
          error: '用户不存在或密码错误',
        },
      };

      mockedAxios.post.mockResolvedValue(mockResponse);

      // Act & Assert
      await expect(authApi.login(loginData)).rejects.toThrow('用户不存在或密码错误');
    });
  });

  describe('validateToken', () => {
    it('should validate token successfully', async () => {
      // Arrange
      const expectedUser = {
        user_id: '123',
        username: 'testuser',
        email: 'test@example.com',
        full_name: 'Test User',
        status: 'active',
      };

      const mockResponse = {
        data: {
          success: true,
          data: expectedUser,
          message: 'Token验证成功',
        },
      };

      mockedAxios.get.mockResolvedValue(mockResponse);

      // Act
      const result = await authApi.validateToken();

      // Assert
      expect(mockedAxios.get).toHaveBeenCalledWith('/auth/validate');
      expect(result).toEqual(expectedUser);
    });

    it('should throw error when token validation fails', async () => {
      // Arrange
      const mockResponse = {
        data: {
          success: false,
          error: 'Token验证失败',
        },
      };

      mockedAxios.get.mockResolvedValue(mockResponse);

      // Act & Assert
      await expect(authApi.validateToken()).rejects.toThrow('Token验证失败');
    });
  });

  describe('logout', () => {
    it('should clear token from localStorage', () => {
      // Act
      authApi.logout();

      // Assert
      expect(localStorageMock.removeItem).toHaveBeenCalledWith('auth_token');
    });
  });
}); 