FROM golang:1.22-alpine AS builder

LABEL stage=gobuilder
# Set proxy environment variables
#ARG PROXY
#RUN export HTTPS_PROXY=$PROXY && export HTTP_PROXY=$PROXY

ENV CGO_ENABLED 0

WORKDIR /build

ADD go.mod .
ADD go.sum .
RUN go mod download
COPY . .
RUN sh ./build.sh

FROM ubuntu:22.04

ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/output /app

CMD ["sh", "./bootstrap.sh"]
