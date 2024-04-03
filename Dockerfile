# 使用 Golang 官方镜像作为基础镜像
FROM golang:alpine as builder

# 设置工作目录
WORKDIR /app

# 将项目文件复制到工作目录
COPY . .

# 编译 Golang 项目
RUN go build -o main .

# 使用轻量级的 alpine 镜像作为基础镜像
FROM alpine

# 设置工作目录
WORKDIR /app

# 将编译好的可执行文件复制到工作目录
COPY --from=builder /app/main .

COPY config.yaml /app/config.yaml

COPY static/ /app/static

COPY Shanghai /usr/share/zoneinfo/Asia/Shanghai
# 暴露应用程序端口
EXPOSE 8080

# 运行应用程序
CMD ["./main"]
