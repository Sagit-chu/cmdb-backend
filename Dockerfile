# 使用官方的 Go 镜像作为构建阶段的基础镜像
FROM golang:1.22-alpine AS builder

# 设置工作目录
WORKDIR /app

# 将 go.mod 和 go.sum 复制到工作目录
COPY go.mod go.sum ./

# 下载所有依赖包。这将被缓存，除非 go.mod 或 go.sum 发生变化
RUN go mod download

# 将项目的源代码复制到工作目录
COPY . .

# 构建应用程序
RUN go build -o cmdb-backend main.go

# 使用一个更小的镜像作为运行阶段的基础镜像
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制二进制文件到当前镜像
COPY --from=builder /app/cmdb-backend .

# 将必要的配置文件（如果有）复制到当前镜像
COPY --from=builder /app/.env .

# 暴露服务端口
EXPOSE 3000

# 运行应用程序
CMD ["./cmdb-backend"]
