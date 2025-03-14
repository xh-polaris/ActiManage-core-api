FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder
# Set proxy environment variables
#ARG PROXY
#RUN export HTTPS_PROXY=$PROXY && export HTTP_PROXY=$PROXY

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
# 设置国内 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct

# 下载依赖
RUN go mod download

COPY . .
RUN sh ./build.sh

FROM ubuntu:22.04

ENV TZ Asia/Tokyo

WORKDIR /app
COPY --from=builder /build/output /app

CMD ["sh", "./bootstrap.sh"]
