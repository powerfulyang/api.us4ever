# Project api.us4ever

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

## 数据库管理

本项目使用 ENT 作为 ORM 框架，使用 Atlas 工具进行数据库结构同步和迁移。

### 配置

数据库配置使用应用程序的统一配置系统，从 Nacos 或环境变量加载。配置结构如下：

```json
{
    "app_name": "api.us4ever",
    "app_env": "development",
    "server": {
        "port": 8080
    },
    "database": {
        "host": "localhost",
        "port": 5432,
        "database": "your_database",
        "username": "your_username",
        "password": "your_password",
        "schema": "public"
    }
}
```

### 命令

以下命令用于管理数据库:

```bash
# 从现有数据库同步结构到 ENT schema
make sync-schema

# 生成数据库迁移
make generate-migration name=迁移名称

# 应用待处理的数据库迁移
make apply-migration

# 初始化数据库连接
make setup-db
```

### 工作流程

1. 同步数据库结构:
   - 从现有数据库读取结构并生成 ENT schema
   - 生成 ENT 代码

2. 迁移管理:
   - 生成迁移文件，记录数据库结构变更
   - 应用迁移文件，更新数据库结构

3. ENT 使用示例:
   ```go
   // 创建数据库服务
   dbService, err := database.New()
   if err != nil {
       log.Fatalf("Failed to create database service: %v", err)
   }

   // 获取 ENT 客户端
   client := dbService.Client()

   // 使用 ENT 客户端操作数据库
   users, err := client.User.Query().All(ctx)
   ```
