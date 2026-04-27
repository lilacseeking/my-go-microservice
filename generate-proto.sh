#!/bin/bash
# scripts/generate-proto.sh
# 用途: 自动生成 proto 对应的 Go 代码

# set -e  # 遇错误立即退出

echo "🔍 正在生成 proto Go 代码..."

protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/user.proto

echo "✅ proto 代码生成完成"
