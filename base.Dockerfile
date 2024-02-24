FROM node:20-alpine as node_base

###

FROM node_base as builder
COPY --from=golang:1.21.5-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add --no-cache binutils file gcc g++ make libc-dev fortify-headers patch git curl wget upx protoc libc6-compat python3 py3-pip && \
    npm i -g pnpm@8 node-gyp

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc docker-entrypoint.sh go.work go.work.sum ./
# generated via
# pnpm gen:dockerfile:copy
# DO NOT CHANGE COMMENTS BELOW

# START COPYGEN
COPY libs/api/package.json libs/api/package.json
COPY libs/brand/package.json libs/brand/package.json
COPY libs/config/package.json libs/config/package.json
COPY libs/crypto/package.json libs/crypto/package.json
COPY libs/fontsource/package.json libs/fontsource/package.json
COPY libs/frontend-chat/package.json libs/frontend-chat/package.json
COPY libs/frontend-now-playing/package.json libs/frontend-now-playing/package.json
COPY libs/grpc/package.json libs/grpc/package.json
COPY libs/pubsub/package.json libs/pubsub/package.json
COPY libs/types/package.json libs/types/package.json
COPY apps/dota/package.json apps/dota/package.json
COPY apps/eval/package.json apps/eval/package.json
COPY apps/integrations/package.json apps/integrations/package.json
COPY frontend/dashboard/package.json frontend/dashboard/package.json
COPY frontend/landing/package.json frontend/landing/package.json
COPY frontend/overlays/package.json frontend/overlays/package.json
COPY frontend/public-page/package.json frontend/public-page/package.json
COPY libs/api/go.mod libs/api/go.mod
COPY libs/config/go.mod libs/config/go.mod
COPY libs/crypto/go.mod libs/crypto/go.mod
COPY libs/gomodels/go.mod libs/gomodels/go.mod
COPY libs/grpc/go.mod libs/grpc/go.mod
COPY libs/helix/go.mod libs/helix/go.mod
COPY libs/integrations/go.mod libs/integrations/go.mod
COPY libs/logger/go.mod libs/logger/go.mod
COPY libs/migrations/go.mod libs/migrations/go.mod
COPY libs/pubsub/go.mod libs/pubsub/go.mod
COPY libs/sentry/go.mod libs/sentry/go.mod
COPY libs/twitch/go.mod libs/twitch/go.mod
COPY libs/types/go.mod libs/types/go.mod
COPY libs/uptrace/go.mod libs/uptrace/go.mod
COPY libs/utils/go.mod libs/utils/go.mod
COPY apps/api/go.mod apps/api/go.mod
COPY apps/bots/go.mod apps/bots/go.mod
COPY apps/discord/go.mod apps/discord/go.mod
COPY apps/emotes-cacher/go.mod apps/emotes-cacher/go.mod
COPY apps/events/go.mod apps/events/go.mod
COPY apps/eventsub/go.mod apps/eventsub/go.mod
COPY apps/parser/go.mod apps/parser/go.mod
COPY apps/scheduler/go.mod apps/scheduler/go.mod
COPY apps/timers/go.mod apps/timers/go.mod
COPY apps/tokens/go.mod apps/tokens/go.mod
COPY apps/websockets/go.mod apps/websockets/go.mod
COPY apps/ytsr/go.mod apps/ytsr/go.mod
# END COPYGEN
COPY libs/helix/go.mod libs/helix/go.mod

# CLI PART
COPY cli cli
COPY libs/config libs/config
COPY libs/crypto libs/crypto
COPY libs/migrations libs/migrations
# CLI PART END

RUN pnpm cli deps

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
