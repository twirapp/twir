FROM node:18-alpine as builder
COPY --from=golang:1.20.1-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add git curl wget upx protoc libc6-compat && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc go.work go.work.sum docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

COPY libs libs
COPY apps apps
COPY frontend frontend
COPY patches patches

RUN pnpm install --frozen-lockfile
RUN pnpm build:libs


### GOLANG MICROSERVICES

FROM alpine:latest as go_prod_base
RUN apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*
COPY --from=builder /app/docker-entrypoint.sh /app/
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM builder as api_builder
RUN cd apps/api && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as api
COPY --from=api_builder /app/apps/api/out /bin/api
CMD ["/bin/api"]

FROM builder as bots_builder
WORKDIR /app
RUN cd apps/bots && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as bots
COPY --from=bots_builder /app/apps/bots/out /bin/bots
CMD ["/bin/bots"]

FROM builder as emotes-cacher_builder
RUN cd apps/emotes-cacher && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as emotes-cacher
COPY --from=emotes-cacher_builder /app/apps/emotes-cacher/out /bin/emotes-cacher
CMD ["/bin/emotes-cacher"]

FROM builder as events_builder
RUN cd apps/events && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as events
COPY --from=events_builder /app/apps/events/out /bin/events
CMD ["/bin/events"]

FROM builder as parser_builder
RUN cd apps/parser && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as parser
COPY --from=parser_builder /app/apps/parser/out /bin/parser
CMD ["/bin/parser"]

FROM builder as streamstatus_builder
RUN cd apps/streamstatus && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as streamstatus
COPY --from=streamstatus_builder /app/apps/api/out /bin/api
CMD ["/bin/api"]

FROM builder as timers_builder
RUN cd apps/timers && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as timers
COPY --from=timers_builder /app/apps/timers/out /bin/timers
CMD ["/bin/timers"]

FROM builder as tokens_builder
RUN cd apps/tokens && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as tokens
COPY --from=tokens_builder /app/apps/tokens/out /bin/tokens
CMD ["/bin/tokens"]

FROM builder as watched_builder
RUN cd apps/watched && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as watched
COPY --from=watched_builder /app/apps/watched/out /bin/watched
CMD ["/bin/watched"]

FROM builder as scheduler_builder
RUN cd apps/scheduler && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as scheduler
COPY --from=scheduler_builder /app/apps/scheduler/out /bin/scheduler
CMD ["/bin/scheduler"]

### NODEJS MICROSERVICES

FROM node:18-alpine as node_prod_base
COPY --from=builder /app/docker-entrypoint.sh /app/
RUN corepack enable
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM builder as dota_builder
RUN cd apps/dota && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as dota
WORKDIR /app
COPY --from=dota_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/dota", "start"]

FROM builder as eval_builder
RUN cd apps/eval && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as eval
WORKDIR /app
COPY --from=eval_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/eval", "start"]

FROM builder as eventsub_builder
RUN cd apps/eventsub && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as eventsub
WORKDIR /app
COPY --from=eventsub_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/eventsub", "start"]

FROM builder as integrations_builder
RUN cd apps/integrations && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as integrations
WORKDIR /app
COPY --from=integrations_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/integrations", "start"]

FROM builder as streamstatus_builder
RUN cd apps/streamstatus && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as streamstatus
WORKDIR /app
COPY --from=streamstatus_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/dota", "start"]

FROM builder as websockets_builder
RUN cd apps/websockets && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as websockets
WORKDIR /app
COPY --from=websockets_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/websockets", "start"]

FROM builder as ytsr_builder
RUN cd apps/ytsr && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as ytsr
WORKDIR /app
COPY --from=ytsr_builder /app /app
CMD ["pnpm", "--filter=@tsuwari/ytsr", "start"]