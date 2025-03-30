# Simple Makefile for a Go project

# Build the application
all: build test

build: generate-ent
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"

# Run import config tool
import-nacos-config:
	@echo "导入配置到Nacos..."
	@go run cmd/nacos-tools/import-config/main.go -file=config/api.us4ever.json

# 打印 Nacos 配置
print-nacos-config:
	@echo "打印 Nacos 配置..."
	@go run ./cmd/nacos-tools/print-config/main.go

.PHONY: all build run test clean watch docker-run docker-down itest build-import-config import-config sync-schema generate-migration apply-migration setup-db print-nacos-config help

# 默认帮助命令
help:
	@echo "数据库管理命令:"
	@echo "  make sync-schema             - 从现有数据库同步结构到 ENT schema"
	@echo ""
	@echo "Nacos 配置命令:"
	@echo "  make import-nacos-config           - 导入配置到 Nacos"
	@echo "  make print-nacos-config      - 打印 Nacos 配置"

# 生成 ENT 代码
generate-ent:
	@echo "生成 ENT 代码..."
	go generate ./ent

# 从数据库同步结构
sync-schema:
	@echo "从数据库同步结构..."
	go run ./cmd/dbtools sync
	@echo "同步完成，生成 ENT 代码..."
	$(MAKE) generate-ent
