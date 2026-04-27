# 测试和数据库初始化脚本说明

## 📁 目录结构

```
scripts/
├── init-db.sh          # Linux/Mac 数据库初始化脚本
├── init-db.ps1         # Windows PowerShell 数据库初始化脚本
├── test-api.sh         # Linux/Mac API 测试脚本
├── test-api.ps1        # Windows PowerShell API 测试脚本
└── README.md           # 本文件
```

## 🗄️ 数据库初始化

### 在 Docker 环境中初始化

如果你使用 Docker Compose 运行服务，可以进入 PostgreSQL 容器执行迁移：

```bash
# 方法 1: 直接进入容器执行
docker exec -i go_micro_postgres psql -U postgres -d myapp_dev < migrations/001_create_users_table.sql

# 方法 2: 使用 docker-compose exec
docker-compose exec postgres psql -U postgres -d myapp_dev -f /docker-entrypoint-initdb.d/001_create_users_table.sql
```

### 在本地环境中初始化

#### Linux/Mac:
```bash
chmod +x scripts/init-db.sh
./scripts/init-db.sh
```

#### Windows PowerShell:
```powershell
.\scripts\init-db.ps1
```

### 手动执行 SQL

```bash
# 连接到 PostgreSQL
psql -h localhost -p 5432 -U postgres -d myapp_dev

# 执行迁移文件
\i migrations/001_create_users_table.sql

# 验证表创建
\dt users

# 查看测试数据
SELECT * FROM users;
```

## 🔧 API 测试

确保应用正在运行（端口 8080），然后执行测试脚本。

### 使用 curl 命令测试（基础示例）

```bash
# 1. 健康检查
curl http://localhost:8081/healthz

# 2. 创建用户
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com"
  }'

# 3. 获取用户（使用返回的 user_id）
curl http://localhost:8080/v1/users/{user_id}

# 4. 获取测试数据中的用户
curl http://localhost:8080/v1/users/test-user-001
```

### 使用自动化测试脚本

#### Linux/Mac:
```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

#### Windows PowerShell:
```powershell
.\scripts\test-api.ps1
```

#### 指定不同的 Base URL:
```bash
# Linux/Mac
./scripts/test-api.sh http://localhost:8080

# Windows
.\scripts\test-api.ps1 -BaseUrl "http://localhost:8080"
```

## 📋 测试用例说明

测试脚本会执行以下测试：

1. **健康检查** - 验证服务是否正常运行
2. **创建用户（成功）** - 测试正常的用户创建流程
3. **创建用户（失败）** - 测试参数验证（空名称、无效邮箱）
4. **获取用户（测试数据）** - 查询初始化时创建的测试用户
5. **获取用户（新创建）** - 查询刚创建的用户
6. **获取用户（不存在）** - 测试 404 错误处理
7. **创建第二个用户** - 测试多次创建

## 🎯 预期响应

### 创建用户成功 (201)
```json
{
  "code": 201,
  "message": "用户创建成功",
  "data": {
    "user_id": "uuid-string"
  }
}
```

### 获取用户成功 (200)
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": "test-user-001",
    "name": "Test User One",
    "email": "test1@example.com",
    "created_at": "2026-04-27T11:00:00Z",
    "updated_at": "2026-04-27T11:00:00Z"
  }
}
```

### 参数验证失败 (400)
```json
{
  "code": 400,
  "message": "参数绑定失败: ..."
}
```

### 用户不存在 (404)
```json
{
  "code": 404,
  "message": "用户不存在"
}
```

## ⚠️ 注意事项

1. **先初始化数据库**：在测试 API 之前，确保已经执行了数据库迁移脚本
2. **Docker 环境**：如果使用 Docker，确保所有容器都已启动且健康
3. **端口冲突**：确保 8080 端口没有被其他应用占用
4. **测试数据**：迁移脚本会插入两条测试数据（test-user-001, test-user-002）

## 🐛 故障排查

### 问题：连接数据库失败
```bash
# 检查 PostgreSQL 是否运行
docker-compose ps postgres

# 查看日志
docker-compose logs postgres
```

### 问题：API 返回 500 错误
```bash
# 检查应用日志
docker-compose logs app

# 确认数据库表已创建
docker exec -it go_micro_postgres psql -U postgres -d myapp_dev -c "\dt users"
```

### 问题：curl 命令不可用
- Windows 10+ 已内置 curl
- 或使用 PowerShell 的 `Invoke-RestMethod`（测试脚本已提供）
- 或安装 Git for Windows（包含 curl）
