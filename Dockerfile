FROM node:16 as builder

WORKDIR /web
COPY ./VERSION .
COPY ./web .

ENV NPM_CONFIG_REGISTRY=https://registry.npmmirror.com/

WORKDIR /web/default
RUN npm install
RUN DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat VERSION) npm run build

WORKDIR /web/berry
RUN npm install
RUN DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat VERSION) npm run build

WORKDIR /web/air
RUN npm install
RUN DISABLE_ESLINT_PLUGIN='true' REACT_APP_VERSION=$(cat VERSION) npm run build

FROM golang:1.20 AS builder2

ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOPROXY="https://goproxy.cn,direct"

#ENV SQL_DSN="oneapi:123456@tcp(127.0.0.1:3300)/one-api"

WORKDIR /build
ADD go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=builder /web/build ./web/build
RUN go build -ldflags "-s -w -X 'github.com/songquanpeng/one-api/common.Version=$(cat VERSION)' -extldflags '-static'" -o one-api

FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates 2>/dev/null || true

COPY --from=builder2 /build/one-api /
EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/one-api"]