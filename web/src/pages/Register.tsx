import React, { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '@/store';
import { authApi } from '@/services/api';
import Button from '@/components/ui/Button';
import Input from '@/components/ui/Input';
import type { CreateUserRequest } from '@/types';

const Register: React.FC = () => {
  const navigate = useNavigate();
  const setAuth = useAuthStore((state) => state.setAuth);
  
  const [formData, setFormData] = useState<CreateUserRequest>({
    username: '',
    email: '',
    password: '',
    full_name: '',
  });
  
  const [errors, setErrors] = useState<Partial<CreateUserRequest>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [apiError, setApiError] = useState<string>('');

  const validateForm = (): boolean => {
    const newErrors: Partial<CreateUserRequest> = {};
    
    if (!formData.username) {
      newErrors.username = '请输入用户名';
    } else if (formData.username.length < 3) {
      newErrors.username = '用户名至少需要3个字符';
    }
    
    if (!formData.email) {
      newErrors.email = '请输入邮箱';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      newErrors.email = '请输入有效的邮箱地址';
    }
    
    if (!formData.password) {
      newErrors.password = '请输入密码';
    } else if (formData.password.length < 6) {
      newErrors.password = '密码至少需要6个字符';
    }
    
    if (!formData.full_name) {
      newErrors.full_name = '请输入全名';
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
      const user = await authApi.register(formData);
      // 注册成功后自动登录
      const loginResponse = await authApi.login({
        email: formData.email,
        password: formData.password,
      });
      setAuth(loginResponse.user, loginResponse.token);
      navigate('/dashboard');
    } catch (error) {
      setApiError(error instanceof Error ? error.message : '注册失败');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    
    // 清除对应字段的错误
    if (errors[name as keyof CreateUserRequest]) {
      setErrors(prev => ({ ...prev, [name]: undefined }));
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            创建新账户
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            已有账户？{' '}
            <Link
              to="/login"
              className="font-medium text-primary-600 hover:text-primary-500"
            >
              立即登录
            </Link>
          </p>
        </div>
        
        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">
            <Input
              label="用户名"
              type="text"
              name="username"
              value={formData.username}
              onChange={handleChange}
              error={errors.username}
              placeholder="请输入用户名"
              required
            />
            
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
              placeholder="请输入密码（至少6个字符）"
              required
            />
            
            <Input
              label="全名"
              type="text"
              name="full_name"
              value={formData.full_name}
              onChange={handleChange}
              error={errors.full_name}
              placeholder="请输入您的全名"
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
            注册
          </Button>
        </form>
      </div>
    </div>
  );
};

export default Register; 