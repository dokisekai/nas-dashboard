# 构建阶段
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 安装构建工具：git（拉依赖）+ gcc/musl-dev（CGO，sqlite 需要）
RUN apk add --no-cache git gcc musl-dev

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用（CGO_ENABLED=1 让 mattn/go-sqlite3 能正常工作）
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o server ./cmd/server

# 运行阶段
FROM alpine:latest

# 安装运行时依赖：ca-certificates（HTTPS）、restic（备份）、rsync（同步）、
# openssh-client（SFTP / 远程 rsync）、rclone（云存储）、tzdata（时区）、
# libc6-compat（musl 兼容层，CGO sqlite 需要）
RUN apk --no-cache add \
      ca-certificates \
      restic \
      rsync \
      openssh-client \
      rclone \
      tzdata \
      libc6-compat \
      shadow \
    && update-ca-certificates

WORKDIR /root/

# 从构建阶段复制二进制、模板、静态前端、种子配置
COPY --from=builder /app/server .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
COPY --from=builder /app/config ./config

# 暴露端口（与 main.go 中保持一致）
EXPOSE 8888

# 让 restic 缓存目录可写
ENV RESTIC_CACHE_DIR=/data/restic-cache \
    TZ=Asia/Shanghai \
    GIN_MODE=release

# 运行应用
CMD ["./server"]
