-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- 插入测试数据（可选）
INSERT INTO users (id, name, email, created_at, updated_at)
VALUES 
    ('test-user-001', 'Test User One', 'test1@example.com', NOW(), NOW()),
    ('test-user-002', 'Test User Two', 'test2@example.com', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;
