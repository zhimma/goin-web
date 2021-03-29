FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    GOPROXY=https://goproxy.io,direct

# 移动到工作目录：/build
WORKDIR /build

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download

# 将代码复制到容器中
COPY ./src .

# 将我们的代码编译成二进制可执行文件 app
RUN go build -o app_server .

###################
# 接下来创建一个小镜像
###################
FROM scratch

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/app_server /
COPY --from=builder /build/config.yaml /

# 需要运行的命令
ENTRYPOINT ["./app_server"]