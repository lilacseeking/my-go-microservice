#!/bin/bash
# 数据库初始化脚本
# 用法: ./scripts/init-db.sh

set -e

echo "========================================="
echo "开始初始化数据库..."
echo "========================================="

# 从 docker-compose 获取数据库连接信息
DB_HOST="${DB_HOST:-localhost}"
DB_PORT="${DB_PORT:-5432}"
DB_USER="${DB_USER:-postgres}"
DB_NAME="${DB_NAME:-myapp_dev}"
DB_PASSWORD="${DB_PASSWORD:-password}"

export PGPASSWORD=$DB_PASSWORD

echo "连接到 PostgreSQL: ${DB_HOST}:${DB_PORT}/${DB_NAME}"

# 等待数据库就绪
echo "等待数据库就绪..."
for i in $(seq 1 30); do
    if pg_isready -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" > /dev/null 2>&1; then
        echo "✓ 数据库已就绪"
        break
    fi
    echo "  等待中... ($i/30)"
    sleep 2
done

# 执行迁移脚本
MIGRATION_DIR="$(dirname "$0")/../migrations"

if [ ! -d "$MIGRATION_DIR" ]; then
    echo "错误: 迁移目录不存在: $MIGRATION_DIR"
    exit 1
fi

echo ""
echo "执行迁移脚本..."
echo "-----------------------------------------"

for migration_file in "$MIGRATION_DIR"/*.sql; do
    if [ -f "$migration_file" ]; then
        filename=$(basename "$migration_file")
        echo "执行: $filename"
        psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$migration_file"
        echo "✓ $filename 执行成功"
        echo ""
    fi
done

echo "========================================="
echo "数据库初始化完成!"
echo "========================================="

# 验证表是否创建成功
echo ""
echo "验证表结构..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "\dt users"
echo ""
echo "查看测试数据..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "SELECT id, name, email FROM users LIMIT 5;"

unset PGPASSWORD
