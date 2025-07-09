# AI开发平台 - 质量工具使用说明

## 概述

为了解决项目规模增大导致的代码质量问题，我们提供了一套完整的解决方案：

## 🛠️ 工具列表

### 1. 测试框架 (`tests/test_framework.go`)
- 统一的测试基础设施
- 数据库事务管理
- 测试数据准备和清理

### 2. 自动化测试脚本 (`scripts/test.sh`)
- 单元测试、集成测试、静态分析
- 代码覆盖率报告
- 测试报告生成

### 3. 模块管理工具 (`scripts/module_manager.sh`)
- 模块依赖分析
- 代码复杂度分析
- 回归测试

### 4. VS Code配置 (`.vscode/tasks.json`)
- 快速运行测试任务
- 代码分析任务
- 开发服务器启动

## 🚀 快速开始

### 1. 环境准备

```bash
# 1. 确保Go环境已安装
go version

# 2. 启动MySQL服务
# macOS: brew services start mysql
# Ubuntu: sudo service mysql start

# 3. 创建测试数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS aicode_test;"

# 4. 安装依赖
go mod tidy
```

### 2. 运行测试

```bash
# 运行所有测试
./scripts/test.sh all

# 只运行单元测试
./scripts/test.sh unit

# 运行静态代码分析
./scripts/test.sh static
```

### 3. 模块管理

```bash
# 查看所有模块
./scripts/module_manager.sh list

# 分析模块复杂度
./scripts/module_manager.sh analyze service

# 检查模块依赖
./scripts/module_manager.sh deps api

# 生成依赖图
./scripts/module_manager.sh graph
```

## 📊 使用VS Code任务

在VS Code中：
1. 按 `Ctrl+Shift+P` (或 `Cmd+Shift+P` on Mac)
2. 输入 "Tasks: Run Task"
3. 选择相应的任务：
   - "运行所有测试"
   - "运行单元测试"
   - "静态代码分析"
   - "模块复杂度分析"

## 📈 质量指标

### 当前项目状态

通过 `./scripts/module_manager.sh analyze service` 分析：

- **Service模块**:
  - 代码行数: 2704
  - 文件数量: 4
  - 函数数量: 81
  - 结构体数量: 17
  - 接口数量: 3
  - 复杂度: 低 (分数: 39)

### 质量标准

- **测试覆盖率**: 目标 > 80%
- **模块复杂度**: 低 (< 50) 或 中 (50-100)
- **函数行数**: < 50行
- **文件行数**: < 500行

## 🔧 开发工作流

### 1. 提交代码前

```bash
# 1. 格式检查
gofmt -l .

# 2. 静态分析
go vet ./...

# 3. 运行单元测试
./scripts/test.sh unit

# 4. 检查复杂度
./scripts/module_manager.sh analyze <module_name>
```

### 2. 重大修改后

```bash
# 运行完整的回归测试
./scripts/module_manager.sh regression
```

### 3. 定期维护

```bash
# 生成质量报告
./scripts/test.sh all

# 查看覆盖率报告
open reports/coverage.html
```

## 📋 最佳实践

### 1. 编写测试

```go
// 好的测试命名
func TestUserService_CreateUser_Success(t *testing.T)
func TestUserService_CreateUser_DuplicateEmail(t *testing.T)

// 测试结构 (AAA模式)
func TestExample(t *testing.T) {
    // Arrange - 准备测试数据
    
    // Act - 执行测试
    
    // Assert - 验证结果
}
```

### 2. 模块设计

- **单一职责**: 每个模块只负责一类功能
- **接口隔离**: 模块间通过接口交互
- **依赖倒置**: 依赖抽象而非具体实现

### 3. 错误处理

```go
func (s *Service) Method() error {
    if err := validate(); err != nil {
        return fmt.Errorf("验证失败: %w", err)
    }
    
    if err := process(); err != nil {
        return fmt.Errorf("处理失败: %w", err)
    }
    
    return nil
}
```

## 🔍 故障排除

### 1. 测试数据库连接失败

```bash
# 检查MySQL服务状态
ps aux | grep mysql

# 启动MySQL服务
brew services start mysql  # macOS
sudo service mysql start   # Ubuntu

# 创建测试数据库
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS aicode_test;"
```

### 2. 依赖安装失败

```bash
# 清理并重新安装依赖
go clean -modcache
go mod tidy
go mod download
```

### 3. 脚本权限问题

```bash
# 添加执行权限
chmod +x scripts/*.sh
```

## 📚 详细文档

更多详细信息请参考：
- [项目质量改进指南](docs/quality_improvement_guide.md)
- [测试策略文档](docs/testing_strategy.md)
- [开发流程文档](development/development_process.md)

## 🤝 贡献指南

1. 每个新功能都需要添加相应的测试
2. 提交代码前运行完整的测试套件
3. 保持代码复杂度在合理范围内
4. 更新相关文档

## 🎯 下一步改进

- [ ] 添加性能测试
- [ ] 集成代码质量扫描工具
- [ ] 自动化CI/CD流程
- [ ] 添加安全性测试
- [ ] 改进测试覆盖率报告

---

通过这套工具，您可以有效地管理项目质量，减少bug，提高开发效率！ 