FROM node:18-alpine as node_base

###

FROM node_base as builder
COPY --from=golang:1.21.0-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add --no-cache build-base git curl wget upx protoc libc6-compat g++ python3

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && \
    npm i -g pnpm@8

COPY . .
RUN chmod +x docker-entrypoint.sh

RUN pnpm install --frozen-lockfile
RUN pnpm turbo run build --filter=./libs/**

###

FROM node_base as node_prod_base
WORKDIR /app
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/*
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
RUN npm i -g pnpm@8
ENTRYPOINT ["/app/docker-entrypoint.sh"]

###

FROM alpine:latest as go_prod_base
WORKDIR /app
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/*
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]
