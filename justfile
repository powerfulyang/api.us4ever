# Windows 上使用 PowerShell 替代 sh:
set windows-shell := ["powershell.exe", "-c"]
set shell := ["sh", "-c"]

# 设置默认命令
default: help

# 构建并测试应用
all: build test

# 构建应用
build:
    @echo "Building..."
    go build -ldflags="-s -w" -o main cmd/api/main.go

# 运行应用
run:
    go run cmd/api/main.go

# 运行测试
test:
    @echo "Testing..."
    go test ./... -v

# 运行集成测试
itest:
    @echo "Running integration tests..."
    go test ./internal/database -v

# 清理二进制文件
clean:
    @echo "Cleaning..."
    rm -f main

# 实时重载
watch:
    #!/usr/bin/env powershell
    if (Get-Command air -ErrorAction SilentlyContinue) {
        air
        Write-Output 'Watching...'
    } else {
        Write-Output 'Installing air...'
        go install github.com/air-verse/air@latest
        air
        Write-Output 'Watching...'
    }

# 导入 Nacos 配置
import-nacos-config:
    @echo "导入配置到Nacos..."
    go run cmd/nacos-tools/import-config/main.go -file=config/api.us4ever.json

# 打印 Nacos 配置
print-nacos-config:
    @echo "打印 Nacos 配置..."
    go run ./cmd/nacos-tools/print-config/main.go

# 生成 ENT 代码
generate-ent:
    @echo "生成 ENT 代码..."
    go generate ./cmd/ent

# 从数据库同步结构
sync-schema: 
    @echo "从数据库同步结构..."
    go run ./cmd/db-tools sync
    @echo "同步完成，生成 ENT 代码..."
    just generate-ent

# 导入数据到数据库
import-moments:
    @echo "导入数据到数据库..."
    go run ./cmd/db-tools import-moments "C:\Users\power\Downloads\simplifyweibo_4_moods.csv"

# 显示帮助信息
help: 
    @echo "数据库管理命令:"
    @echo "  just sync-schema             - 从现有数据库同步结构到 ENT schema"
    @echo ""
    @echo "Nacos 配置命令:"
    @echo "  just import-nacos-config     - 导入配置到 Nacos"
    @echo "  just print-nacos-config      - 打印 Nacos 配置"