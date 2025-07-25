/**
 * 认证相关的工具函数
 */

/**
 * 获取认证token
 */
export const getAuthToken = (): string | null => {
  return localStorage.getItem('auth_token');
};

/**
 * 获取认证头
 */
export const getAuthHeaders = () => {
  const token = getAuthToken();
  return token ? { 'Authorization': `Bearer ${token}` } : {};
};

/**
 * 设置认证token
 */
export const setAuthToken = (token: string) => {
  localStorage.setItem('auth_token', token);
};

/**
 * 清除认证token
 */
export const clearAuthToken = () => {
  localStorage.removeItem('auth_token');
  localStorage.removeItem('auth-store');
};