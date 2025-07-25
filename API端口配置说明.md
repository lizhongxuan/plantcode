# API端口配置说明

## 端口架构

### 开发环境
- **前端开发服务器**: `http://localhost:3000` (Vite Dev Server)
- **后端API服务器**: `http://localhost:8080` (Go Server)

### API代理配置

前端使用 Vite 开发服务器的代理功能，将API请求自动转发到后端服务器：

```typescript
// vite.config.ts
server: {
  port: 3000,
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    },
    '/health': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    }
  }
}
```

## 请求流程

1. **前端页面访问**: `http://localhost:3000`
2. **JavaScript API调用**: `/api/projects/list`
3. **Vite代理转发**: `http://localhost:3000/api/projects/list` → `http://localhost:8080/api/projects/list`
4. **后端处理**: Go服务器在8080端口处理请求
5. **响应返回**: 后端响应通过代理返回给前端

## 正确的访问方式

### ✅ 正确方式
- 访问前端页面: `http://localhost:3000`
- 在前端页面中，JavaScript会自动通过代理调用API

### ❌ 错误方式
- 直接在浏览器访问: `http://localhost:3000/api/projects/list`
  - 这会被代理转发到后端，但没有前端页面的认证token，会返回401错误

## API配置

前端API配置文件位于 `web/src/config/api.ts`:

```typescript
const config = {
  development: {
    baseURL: '/api',  // 使用相对路径，通过Vite代理转发
    timeout: 10000,
  },
  production: {
    baseURL: '/api',  // 生产环境下前后端在同一域名
    timeout: 10000,
  }
};
```

## 启动顺序

1. **启动后端服务器** (8080端口):
   ```bash
   go run cmd/server/main.go
   ```

2. **启动前端开发服务器** (3000端口):
   ```bash
   cd web && npm run dev
   ```

3. **访问应用**: `http://localhost:3000`

## 故障排查

### 如果看到请求到3000端口的API错误:

1. **检查前端开发服务器是否运行**:
   ```bash
   lsof -i :3000
   ```

2. **检查后端服务器是否运行**:
   ```bash
   lsof -i :8080
   ```

3. **验证代理配置**:
   查看 `web/vite.config.ts` 中的 proxy 配置

4. **清除缓存重启**:
   ```bash
   cd web
   rm -rf node_modules/.vite
   npm run dev
   ```

### 如果直接访问API端点:

- 应该访问: `http://localhost:8080/api/projects/list` (带认证token)
- 而不是: `http://localhost:3000/api/projects/list`

## 生产环境

在生产环境中，通常前端静态文件和后端API会部署在同一域名下，所以使用相对路径 `/api` 是正确的。