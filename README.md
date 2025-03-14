# go-ext

这是一个工具集合，包含配置、日志记录、RESTful 及其他工具，帮助快速开发 http 相关应用程序。

## 特性

- **配置管理**
  - 支持多种数据库类型（MySQL、SQLite3）
  - 灵活的数据库连接配置选项
  - 使用mapstructure进行结构化配置
  - 支持yaml、json、toml配置文件和环境变量（EXT_xxx）

- **日志系统**
  - 多种日志级别（Debug、Info、Warn、Error、Fatal）
  - 可配置的输出格式和目标
  - 与 zap logger 集成

- **RESTful API框架**
  - 基于Gin框架构建
  - 标准化的REST端点（GET、POST、PUT、DELETE、PATCH）
  - 中间件支持

- **存储层**
  - 数据库抽象
  - 连接池
  - 查询构建器

## 安装

```bash
go get github.com/fize/go-ext
```

## 快速开始
### 数据库配置
```go
import "github.com/fize/go-ext/config"

// Create database configuration
sqlConfig, err := config.NewSQLConfig(
    config.WithType("mysql"),
    config.WithHost("localhost:3306"),
    config.WithUser("root"),
    config.WithPassword("password"),
    config.WithDB("myapp"),
)
```
### 日志配置
```go
import "github.com/fize/go-ext/log"

// Use global logger
log.Info("Starting application...")
log.Debugf("Connected to database: %s", dbName)
log.Error("Failed to process request", err)
```