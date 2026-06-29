# 构建阶段：编译 Go 后端
FROM golang:1.24 AS build
WORKDIR /src
# 先单独复制依赖清单，利用分层缓存：依赖不变时无需重新下载
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# 纯 Go 构建（sqlite/mysql 驱动均为纯 Go 实现），产出静态可执行文件
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# 运行阶段：精简运行镜像，只保留可执行文件
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/server /app/server
# 视频上传目录，可通过数据卷挂载持久化
RUN mkdir -p /app/uploads
ENV APP_PORT=8080 \
    UPLOAD_DIR=/app/uploads \
    TZ=Asia/Shanghai
EXPOSE 8080
ENTRYPOINT ["/app/server"]
