FROM node:18-alpine as node_base

###

FROM node_base as builder
COPY --from=golang:1.21.0-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add --no-cache build-base git curl wget upx protoc libc6-compat python3

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && \
    npm i -g pnpm@8

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc docker-entrypoint.sh ./
# generated via
# pnpm gen:dockerfile:packagecopy
COPY libs/config/package.json libs/config/package.json
COPY libs/crypto/package.json libs/crypto/package.json
COPY libs/frontend-chat/package.json libs/frontend-chat/package.json
COPY libs/grpc/package.json libs/grpc/package.json
COPY libs/pubsub/package.json libs/pubsub/package.json
COPY libs/types/package.json libs/types/package.json
COPY apps/dota/package.json apps/dota/package.json
COPY apps/eval/package.json apps/eval/package.json
COPY apps/integrations/package.json apps/integrations/package.json
COPY apps/language-detector/package.json apps/language-detector/package.json
COPY frontend/dashboard/package.json frontend/dashboard/package.json
COPY frontend/landing/package.json frontend/landing/package.json
COPY frontend/overlays/package.json frontend/overlays/package.json
COPY frontend/public-page/package.json frontend/public-page/package.json
RUN pnpm fetch
RUN chmod +x docker-entrypoint.sh

COPY . .


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
