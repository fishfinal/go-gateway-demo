# Go Gateway Demo

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Gin Framework](https://img.shields.io/badge/Gin-1.9+-00ADD8?style=flat&logo=gin)](https://gin-gonic.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

> 分布式高可用入口架构实战系列示例库

这是一个使用 **Go语言** 和 **Gin框架** 实现的轻量级 Web 服务示例，专为《分布式高可用入口架构实战系列》教程设计。它核心实现了健康检查与身份标识功能，是后续搭建 **Nginx + Keepalived** 高可用集群的基础后端服务。

## ✨ 核心功能

本项目提供了一个“会说话”的 Web 服务，包含以下接口，用于验证和模拟高可用架构中的各种场景：

| 功能 | 接口 | 用途 |
| :--- | :--- | :--- |
| **健康检查** | `GET /health` | 供 Nginx/Keepalived 探测，返回服务状态与身份信息。 |
| **业务模拟** | `POST /api/v1/task` | 模拟任务下发，返回处理请求的节点信息。 |
| **慢接口模拟** | `GET /api/slow` | 模拟 5 秒延迟响应，用于测试超时和熔断机制。 |
| **错误接口模拟** | `GET /api/error` | 返回 500 错误，用于测试故障摘除和节点下线。 |
| **调试信息** | `GET /info` | 查看服务所在主机的 Hostname、IP 和端口。 |

## 🚀 快速开始

你可以通过以下三种方式之一快速运行此服务。

### 1. 直接运行 (本地开发)

```bash
# 克隆项目
git clone https://github.com/fishfinal/go-gateway-demo.git
cd go-gateway-demo

# 下载依赖
go mod download

# 运行服务 (默认监听 8080 端口)
go run main.go
```

### 2. 使用 Docker 运行

```bash
# 构建镜像
docker build -t go-gateway-demo .

# 运行容器
docker run -d -p 8080:8080 --name gateway-demo go-gateway-demo
```

### 3. 使用 Docker Compose 模拟三节点集群

这是本教程推荐的方式，用于在一台机器上模拟三台服务器的集群环境。

```bash
# 启动三个节点，分别映射到宿主机的 8081, 8082, 8083 端口
docker-compose up -d --build

# 查看集群状态
docker-compose ps
```

## 📡 API 测试示例

服务启动后，你可以使用 `curl` 命令进行测试。

### 健康检查

```bash
curl http://localhost:8080/health | jq .
```

**响应示例**:
```json
{
  "status": "ok",
  "hostname": "your-hostname",
  "ip": "192.168.x.x",
  "timestamp": 1749876543
}
```

### 业务接口

```bash
curl -X POST http://localhost:8080/api/v1/task \
  -H "Content-Type: application/json" \
  -d '{"task_id": "task-001", "target": "192.168.0.0/24"}' | jq .
```

**响应示例**:
```json
{
  "code": 0,
  "message": "task accepted",
  "processed_by": {
    "hostname": "your-hostname",
    "ip": "192.168.x.x"
  },
  "task": {
    "task_id": "task-001",
    "target": "192.168.0.0/24"
  }
}
```

## 📦 部署到生产环境

对于生产环境，建议使用 `systemd` 进行服务管理。

1.  **交叉编译 Linux 二进制文件**:
    ```bash
    GOOS=linux GOARCH=amd64 go build -o gateway main.go
    ```

2.  **上传并配置 `systemd` 服务**:
    将生成的 `gateway` 二进制文件上传至服务器，并参考以下配置创建服务文件。
    ```bash
    # /etc/systemd/system/gateway.service
    [Unit]
    Description=Gateway Health Check Service
    After=network.target

    [Service]
    Type=simple
    User=www-data
    WorkingDirectory=/var/lib/gateway/
    ExecStart=/usr/local/bin/gateway
    Restart=always
    RestartSec=5
    Environment="PORT=8080"

    [Install]
    WantedBy=multi-user.target
    ```

## 📂 项目结构

```
go-gateway-demo/
├── main.go                 # 主程序入口，包含所有路由和处理函数
├── go.mod                  # Go 模块依赖管理
├── go.sum                  # Go 模块依赖校验
├── Dockerfile              # Docker 镜像构建文件
├── docker-compose.yml      # 三节点本地集群编排文件
└── README.md               # 项目说明文档
```

## 🔗 相关教程

本项目是 **《分布式高可用入口架构实战系列》** 的配套代码，你可以通过以下文章了解更详细的架构原理和部署步骤：

- [第1篇：高可用入口架构概览](https://fishfinal.com/posts/high-availability-gateway/01-architecture-overview.html)
- [第2篇：VRRP 协议原理与 Keepalived 基础](https://fishfinal.com/posts/high-availability-gateway/02-vrrp-keepalived.html)
- **[第3篇：Golang + Gin 实现健康检查服务 (本仓库)](https://fishfinal.com/posts/high-availability-gateway/03-golang-gin-healthcheck-service.html)**
- [第4篇：Nginx 健康检查与主动熔断 (待发布)](https://fishfinal.com/posts/high-availability-gateway/04-nginx-healthcheck.html)

## 🤝 贡献与反馈

欢迎通过 Issue 或 Pull Request 提出建议和改进。如果你觉得这个项目有帮助，请给它一个 ⭐️！
