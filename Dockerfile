FROM node:18-alpine as builder
COPY --from=golang:1.20.1-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app


RUN apk add git curl wget upx protoc libc6-compat && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    npm i -g pnpm@7

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc go.work go.work.sum docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

COPY libs libs
COPY apps apps
COPY patches patches

RUN pnpm install --frozen-lockfile
RUN pnpm build:libs


FROM alpine:latest as go_prod_base
RUN apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*
COPY --from=builder /app/docker-entrypoint.sh /app/
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM golang:1.20.1-alpine as api_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/api && go mod download
RUN cd apps/api && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as api
COPY --from=api_builder /app/apps/api/out /bin/api
CMD ["/bin/api"]

FROM golang:1.20.1-alpine as bots_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/bots && go mod download
RUN cd apps/bots && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as bots
COPY --from=bots_builder /app/apps/bots/out /bin/bots
CMD ["/bin/bots"]

FROM golang:1.20.1-alpine as emotes-cacher_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/emotes-cacher && go mod download
RUN cd apps/emotes-cacher && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as emotes-cacher
COPY --from=emotes-cacher_builder /app/apps/emotes-cacher/out /bin/emotes-cacher
CMD ["/bin/emotes-cacher"]

FROM golang:1.20.1-alpine as events_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/events && go mod download
RUN cd apps/events && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as events
COPY --from=events_builder /app/apps/events/out /bin/events
CMD ["/bin/events"]

FROM golang:1.20.1-alpine as parser_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/parser && go mod download
RUN cd apps/parser && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as parser
COPY --from=parser_builder /app/apps/parser/out /bin/parser
CMD ["/bin/parser"]

FROM golang:1.20.1-alpine as streamstatus_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/streamstatus && go mod download
RUN cd apps/streamstatus && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as streamstatus
COPY --from=streamstatus_builder /app/apps/api/out /bin/api
CMD ["/bin/api"]

FROM golang:1.20.1-alpine as timers_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/timers && go mod download
RUN cd apps/timers && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as timers
COPY --from=timers_builder /app/apps/timers/out /bin/timers
CMD ["/bin/timers"]

FROM golang:1.20.1-alpine as tokens_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/tokens && go mod download
RUN cd apps/tokens && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as tokens
COPY --from=tokens_builder /app/apps/tokens/out /bin/tokens
CMD ["/bin/tokens"]

FROM golang:1.20.1-alpine as watched_builder
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=builder /app/apps /app/apps
COPY --from=builder /app/libs /app/libs
RUN cd apps/watched && go mod download
RUN cd apps/watched && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as watched
COPY --from=watched_builder /app/apps/watched/out /bin/watched
CMD ["/bin/watched"]
