FROM node:18-alpine as builder
COPY --from=golang:1.21.0-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

WORKDIR /app

RUN apk add git curl wget upx protoc libc6-compat g++ && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0 && \
    go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && \
    npm i -g pnpm@8

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc go.work go.work.sum docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

COPY libs libs
COPY tools tools
COPY apps apps
COPY frontend frontend
COPY patches patches

RUN pnpm install --frozen-lockfile
RUN pnpm build:libs

### GOLANG MICROSERVICES

FROM alpine:latest as go_prod_base
RUN  apk add wget && \
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

FROM builder as ytsr_builder
WORKDIR /app
RUN cd apps/ytsr && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as ytsr
COPY --from=ytsr_builder /app/apps/ytsr/out /bin/ytsr
CMD ["/bin/ytsr"]

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

FROM builder as scheduler_builder
RUN cd apps/scheduler && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as scheduler
COPY --from=scheduler_builder /app/apps/scheduler/out /bin/scheduler
CMD ["/bin/scheduler"]

FROM builder as websockets_builder
RUN cd apps/websockets && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as websockets
COPY --from=websockets_builder /app/apps/websockets/out /bin/websockets
CMD ["/bin/websockets"]

FROM builder as eventsub_builder
RUN cd apps/eventsub && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as eventsub
COPY --from=eventsub_builder /app/apps/eventsub/out /bin/eventsub
CMD ["/bin/eventsub"]

FROM builder as discord_builder
RUN cd apps/discord && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as discord
COPY --from=discord_builder /app/apps/discord/out /bin/discord
CMD ["/bin/discord"]

### NODEJS MICROSERVICES

FROM node:18-alpine as node_prod_base
WORKDIR /app
RUN apk add wget g++ && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/*
COPY --from=builder /app/package.json /app/pnpm-lock.yaml /app/pnpm-workspace.yaml /app/.npmrc /app/docker-entrypoint.sh ./
COPY --from=builder /app/node_modules/ ./node_modules/
RUN chmod +x docker-entrypoint.sh
RUN npm i -g pnpm@8
ENTRYPOINT ["/app/docker-entrypoint.sh"]

#FROM builder as dota_builder
#RUN cd apps/dota && \
#    pnpm build && \
#    pnpm prune --prod
#
#FROM node_prod_base as dota
#WORKDIR /app
#COPY --from=dota_builder /app/apps/dota /app/apps/dota
#COPY --from=dota_builder /app/libs/config /app/libs/config
#COPY --from=dota_builder /app/libs/grpc /app/libs/grpc
#COPY --from=dota_builder /app/libs/shared /app/libs/shared
#COPY --from=dota_builder /app/libs/typeorm /app/libs/typeorm
#CMD ["pnpm", "--filter=@twir/dota", "start"]

FROM builder as eval_builder
RUN cd apps/eval && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as eval
WORKDIR /app
COPY --from=eval_builder /app/apps/eval /app/apps/eval
COPY --from=eval_builder /app/libs/config /app/libs/config
COPY --from=eval_builder /app/libs/grpc /app/libs/grpc
CMD ["pnpm", "--filter=@twir/eval", "start"]

FROM builder as language-detector_builder
RUN cd apps/language-detector && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as language-detector
WORKDIR /app
COPY --from=language-detector_builder /app/apps/language-detector /app/apps/language-detector
COPY --from=language-detector_builder /app/libs/grpc /app/libs/grpc
CMD ["pnpm", "--filter=@twir/language-detector", "start"]

FROM builder as integrations_builder
RUN cd apps/integrations && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as integrations
WORKDIR /app
COPY --from=integrations_builder /app/apps/integrations /app/apps/integrations
COPY --from=integrations_builder /app/libs/config /app/libs/config
COPY --from=integrations_builder /app/libs/grpc /app/libs/grpc
COPY --from=integrations_builder /app/libs/pubsub /app/libs/pubsub
CMD ["pnpm", "--filter=@twir/integrations", "start"]

### FRONTEND

FROM builder as dashboard_builder
RUN cd frontend/dashboard && \
    pnpm build

FROM caddy:latest as dashboard
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=dashboard_builder /app/frontend/dashboard/dist/ /app
EXPOSE 80

FROM builder as public-page_builder
RUN cd frontend/public-page && \
    pnpm build

FROM caddy:latest as public
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=public-page_builder /app/frontend/public-page/dist/ /app
EXPOSE 80

FROM builder as landing_builder
RUN cd frontend/landing && \
    pnpm build && \
    pnpm prune --prod

FROM node_prod_base as landing
WORKDIR /app
COPY --from=landing_builder /app/frontend/landing /app/frontend/landing
COPY --from=landing_builder /app/libs/config /app/libs/config
COPY --from=landing_builder /app/libs/grpc /app/libs/grpc
CMD ["pnpm", "--filter=@twir/landing", "start"]

FROM builder as overlays_builder
RUN cd frontend/overlays && \
    pnpm build

FROM steebchen/nginx-spa:stable as overlays
COPY --from=overlays_builder /app/frontend/overlays/dist/ /app
EXPOSE 80
CMD ["nginx"]
### MIGRATIONS

FROM builder as migrations_builder
RUN cd libs/migrations && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./main.go && upx -9 -k ./out

FROM go_prod_base as migrations
COPY --from=migrations_builder /app/libs/migrations/out /bin/migrations
CMD ["sh", "-c", "/bin/migrations && sleep 60"]
