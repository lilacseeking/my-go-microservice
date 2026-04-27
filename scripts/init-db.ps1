# 数据库初始化脚本 (PowerShell版本)
# 用法: .\scripts\init-db.ps1

$ErrorActionPreference = "Stop"

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "开始初始化数据库..." -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 数据库连接信息
$DB_HOST = if ($env:DB_HOST) { $env:DB_HOST } else { "localhost" }
$DB_PORT = if ($env:DB_PORT) { $env:DB_PORT } else { "5432" }
$DB_USER = if ($env:DB_USER) { $env:DB_USER } else { "postgres" }
$DB_NAME = if ($env:DB_NAME) { $env:DB_NAME } else { "myapp_dev" }
$DB_PASSWORD = if ($env:DB_PASSWORD) { $env:DB_PASSWORD } else { "password" }

$env:PGPASSWORD = $DB_PASSWORD

Write-Host "连接到 PostgreSQL: ${DB_HOST}:${DB_PORT}/${DB_NAME}" -ForegroundColor Yellow

# 等待数据库就绪
Write-Host "等待数据库就绪..." -ForegroundColor Yellow
for ($i = 1; $i -le 30; $i++) {
    try {
        $result = psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT 1;" 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ 数据库已就绪" -ForegroundColor Green
            break
        }
    } catch {
        # 继续等待
    }
    Write-Host "  等待中... ($i/30)" -ForegroundColor Gray
    Start-Sleep -Seconds 2
}

# 获取迁移目录
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$MigrationDir = Join-Path $ScriptDir "..\migrations"

if (-not (Test-Path $MigrationDir)) {
    Write-Host "错误: 迁移目录不存在: $MigrationDir" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "执行迁移脚本..." -ForegroundColor Yellow
Write-Host "-----------------------------------------" -ForegroundColor Gray

# 执行所有SQL迁移文件
Get-ChildItem -Path $MigrationDir -Filter "*.sql" | Sort-Object Name | ForEach-Object {
    $migrationFile = $_.FullName
    $filename = $_.Name
    Write-Host "执行: $filename" -ForegroundColor Cyan
    
    try {
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migrationFile
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✓ $filename 执行成功" -ForegroundColor Green
        } else {
            Write-Host "✗ $filename 执行失败" -ForegroundColor Red
            exit 1
        }
    } catch {
        Write-Host "✗ $filename 执行出错: $_" -ForegroundColor Red
        exit 1
    }
    Write-Host ""
}

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host "数据库初始化完成!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan

# 验证表是否创建成功
Write-Host ""
Write-Host "验证表结构..." -ForegroundColor Yellow
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\dt users"

Write-Host ""
Write-Host "查看测试数据..." -ForegroundColor Yellow
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "SELECT id, name, email FROM users LIMIT 5;"

Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
