# API 测试命令 (CURL Examples)

## 基础信息

- **HTTP API 地址**: `http://localhost:8080`
- **健康检查地址**: `http://localhost:8081/healthz`
- **API 版本**: v1

## 快速开始

### 1. 初始化数据库

在使用API之前，需要先执行数据库迁移：

```powershell
# Windows PowerShell
docker cp migrations/001_create_users_table.sql go_micro_postgres:/tmp/init.sql
docker exec go_micro_postgres psql -U postgres -d myapp_dev -f /tmp/init.sql
```

```bash
# Linux/Mac
docker exec -i go_micro_postgres psql -U postgres -d myapp_dev < migrations/001_create_users_table.sql
```

### 2. 测试API

#### 健康检查

```bash
curl http://localhost:8081/healthz
```

**预期响应**: `OK` (HTTP 200)

---

#### 创建用户

```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com"
  }'
```

**预期响应** (HTTP 201):
```json
{
  "code": 201,
  "message": "用户创建成功",
  "data": {
    "user_id": "uuid-string-here"
  }
}
```

---

#### 获取用户

使用上面返回的 `user_id`：

```bash
curl http://localhost:8080/v1/users/{user_id}
```

或者使用测试数据中的ID：

```bash
curl http://localhost:8080/v1/users/test-user-001
```

**预期响应** (HTTP 200):
```json
{
  "code": 200,
  "message": "查询成功",
  "data": {
    "id": "test-user-001",
    "name": "Test User One",
    "email": "test1@example.com",
    "created_at": "2026-04-27T11:37:45.848482Z",
    "updated_at": "2026-04-27T11:37:45.848482Z"
  }
}
```

---

#### 错误场景测试

**参数验证失败** (HTTP 400):
```bash
curl -X POST http://localhost:8080/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "",
    "email": "invalid-email"
  }'
```

**用户不存在** (HTTP 200, data为null):
```bash
curl http://localhost:8080/v1/users/non-existent-id
```

---

## PowerShell 示例

如果你使用的是Windows PowerShell，可以使用以下命令：

### 创建用户
```powershell
$body = @{
    name = "张三"
    email = "zhangsan@example.com"
} | ConvertTo-Json

Invoke-RestMethod -Uri 'http://localhost:8080/v1/users' `
    -Method Post `
    -ContentType 'application/json' `
    -Body $body
```

### 获取用户
```powershell
Invoke-RestMethod -Uri 'http://localhost:8080/v1/users/test-user-001' -Method Get
```

### 健康检查
```powershell
Invoke-RestMethod -Uri 'http://localhost:8081/healthz' -Method Get
```

---

## 自动化测试

项目提供了自动化测试脚本：

### Linux/Mac
```bash
chmod +x scripts/test-api.sh
./scripts/test-api.sh
```

### Windows PowerShell
```powershell
.\scripts\test-api.ps1
```

测试脚本会自动执行7个测试用例并生成测试报告。

---

## 常见问题

### Q: 创建用户时返回500错误
A: 可能是邮箱已存在。尝试使用不同的邮箱地址。

### Q: 获取用户时返回200但data为null
A: 用户ID不存在。确认ID是否正确，或先创建用户。

### Q: 连接被拒绝
A: 确保Docker容器正在运行：
```bash
docker-compose ps
```

### Q: 如何查看应用日志
```bash
docker logs go_micro_app --tail 50
```

---

## API 端点总结

| 方法 | 路径 | 描述 | 状态码 |
|------|------|------|--------|
| GET | `/healthz` | 健康检查（8081端口） | 200 |
| POST | `/v1/users` | 创建用户 | 201, 400, 500 |
| GET | `/v1/users/:id` | 获取用户 | 200, 404 |

---

## 测试数据

数据库初始化后会创建以下测试用户：

| ID | 名称 | 邮箱 |
|----|------|------|
| test-user-001 | Test User One | test1@example.com |
| test-user-002 | Test User Two | test2@example.com |
