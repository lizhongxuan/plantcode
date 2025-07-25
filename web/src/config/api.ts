// API配置
const config = {
  // 开发环境使用相对路径，生产环境可以使用绝对路径
  development: {
    baseURL: '/api',  // 通过Vite代理转发到8080端口
    timeout: 10000,
  },
  production: {
    baseURL: '/api',  // 生产环境下后端和前端在同一域名
    timeout: 10000,
  }
};

// 获取当前环境
const isDevelopment = import.meta.env.DEV;
const environment = isDevelopment ? 'development' : 'production';

// 导出当前环境的配置
export const apiConfig = config[environment];

// 导出后端服务器地址（仅用于开发调试）
export const BACKEND_URL = 'http://localhost:8080';

// 导出前端开发服务器地址
export const FRONTEND_URL = 'http://localhost:3000';

// 环境信息
export const ENV_INFO = {
  isDevelopment,
  environment,
  frontendPort: 3000,
  backendPort: 8080,
};