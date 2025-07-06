import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { authApi } from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import type { LoginRequest } from '@/types';

const Login: React.FC = () => {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  
  const [formData, setFormData] = useState<LoginRequest>({
    email: '',
    password: '',
  });
  
  const [errors, setErrors] = useState<Partial<LoginRequest>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [apiError, setApiError] = useState<string>('');

  const validateForm = (): boolean => {
    const newErrors: Partial<LoginRequest> = {};
    
    if (!formData.email) {
      newErrors.email = '请输入邮箱';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = '请输入有效的邮箱地址';
    }
    
    if (!formData.password) {
      newErrors.password = '请输入密码';
    }
    
    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setApiError('');
    
    if (!validateForm()) return;
    
    setIsLoading(true);
    
    try {
      const response = await authApi.login(formData);
      setAuth(response.user, response.token);
      navigate('/dashboard');
    } catch (error) {
      setApiError(error instanceof Error ? error.message : '登录失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // 清除对应字段的错误
    if (errors[name as keyof LoginRequest]) {
      setErrors(prev => ({ ...prev, [name]: undefined }));
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            登录到您的账户
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            还没有账户？{' '}
            <Link
              to="/register"
              className="font-medium text-primary-600 hover:text-primary-500"
            >
              立即注册
            </Link>
          </p>
        </div>
        
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <Input
              label="邮箱地址"
              type="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              error={errors.email}
              placeholder="请输入您的邮箱"
              required
            />
            
            <Input
              label="密码"
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              error={errors.password}
              placeholder="请输入您的密码"
              required
            />
          </div>

          {apiError && (
            <div className="rounded-md bg-red-50 p-4">
              <p className="text-sm text-red-800">{apiError}</p>
            </div>
          )}

          <Button
            type="submit"
            isLoading={isLoading}
            className="w-full"
            size="lg"
          >
            登录
          </Button>
        </form>
      </div>
    </div>
  );
};

export default Login; 