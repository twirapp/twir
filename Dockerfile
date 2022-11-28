FROM node:18-alpine as base

COPY --from=golang:1.19.2-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
RUN apk add --no-cache protoc git curl

WORKDIR /app
RUN npm i -g pnpm@7

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc go.mod go.work go.work.sum docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh

COPY libs libs
COPY apps apps

RUN pnpm install --filter=!ngrok
RUN pnpm build

FROM node:18-alpine as node_prod_base
RUN npm i -g pnpm
RUN apk add wget
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub
RUN echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories
RUN apk add doppler
RUN rm -rf /var/cache/apk/*

FROM node:18-alpine as node_deps_base
WORKDIR /app
RUN npm i -g pnpm
RUN apk add git
COPY --from=base /app/package.json /app/pnpm-lock.yaml /app/pnpm-workspace.yaml /app/turbo.json /app/.npmrc ./

FROM node_deps_base as dota_deps
RUN apk add openssh
COPY --from=base /app/apps/dota apps/dota/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/shared libs/shared/
COPY --from=base /app/libs/nats libs/nats/
RUN pnpm install --prod

FROM node_prod_base as dota
WORKDIR /app
COPY --from=dota_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:dota"]

FROM node_deps_base as eval_deps
COPY --from=base /app/apps/eval apps/eval/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/nats libs/nats/
RUN pnpm install --prod

FROM node_prod_base as eval
WORKDIR /app
COPY --from=eval_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:eval"]

FROM node_deps_base as eventsub_deps
COPY --from=base /app/apps/eventsub apps/eventsub/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/nats libs/nats/
COPY --from=base /app/libs/shared libs/shared/
COPY --from=base /app/libs/typeorm libs/typeorm/
RUN pnpm install --prod

FROM node_prod_base as eventsub
WORKDIR /app
COPY --from=eventsub_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:eventsub"]

FROM node_deps_base as integrations_deps
COPY --from=base /app/apps/integrations apps/integrations/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/nats libs/nats/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/shared libs/shared/
RUN pnpm install --prod

FROM node_prod_base as integrations
WORKDIR /app
COPY --from=integrations_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:integrations"]

FROM node_deps_base as scheduler_deps
COPY --from=base /app/apps/scheduler apps/scheduler/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/nats libs/nats/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/shared libs/shared/
RUN pnpm install --prod

FROM node_prod_base as scheduler
WORKDIR /app
COPY --from=scheduler_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:scheduler"]

FROM node_deps_base as streamstatus_deps
COPY --from=base /app/apps/streamstatus apps/streamstatus/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/shared libs/shared/
COPY --from=base /app/libs/nats libs/nats/
RUN pnpm install --prod

FROM node_prod_base as streamstatus
WORKDIR /app
COPY --from=streamstatus_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:streamstatus"]

FROM node_deps_base as migrations_deps
COPY --from=base /app/tsconfig.json /app/tsconfig.base.json ./
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/config libs/config/
RUN pnpm install --prod

FROM node_prod_base as migrations
WORKDIR /app
COPY --from=migrations_deps /app/ /app/
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "run", "migrate:deploy"]

FROM node_deps_base as web_deps
COPY --from=base /app/apps/web apps/web/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/shared libs/shared/
COPY --from=base /app/libs/ui libs/ui/
COPY --from=base /app/libs/config libs/config/
RUN pnpm install --prod

FROM node_prod_base as web
WORKDIR /app
COPY --from=web_deps /app/ /app/
EXPOSE 3000
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["pnpm", "start:web"]

FROM node_deps_base as dashboard_deps
COPY --from=base /app/apps/frontend apps/frontend/
COPY --from=base /app/libs/shared libs/shared/
COPY --from=base /app/libs/typeorm libs/typeorm/
COPY --from=base /app/libs/config libs/config/
RUN pnpm install --prod

FROM devforth/spa-to-http:latest as dashboard
WORKDIR /app
COPY --from=dashboard_deps /app/apps/frontend/dist .
EXPOSE 8080

FROM alpine:latest as go_prod_base
RUN apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*

FROM golang:1.19.2-alpine as golang_deps_base
WORKDIR /app
RUN apk add git curl wget upx
COPY --from=base /app/apps/parser apps/parser/
COPY --from=base /app/apps/timers apps/timers/
COPY --from=base /app/apps/bots apps/bots/
COPY --from=base /app/apps/api apps/api/
COPY --from=base /app/apps/watched apps/watched/
COPY --from=base /app/libs/config libs/config/
COPY --from=base /app/libs/nats libs/nats/
COPY --from=base /app/libs/twitch libs/twitch/
COPY --from=base /app/libs/gomodels libs/gomodels/
COPY --from=base /app/libs/integrations/spotify libs/integrations/spotify/
RUN rm -r `find . -name node_modules -type d`

FROM golang_deps_base as parser_deps
RUN cd apps/parser && go mod download
RUN cd apps/parser && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as parser
COPY --from=parser_deps /app/apps/parser/out /bin/parser
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/parser"]

FROM golang_deps_base as timers_deps
RUN cd apps/timers && go mod download
RUN cd apps/timers && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as timers
COPY --from=timers_deps /app/apps/timers/out /bin/timers
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/timers"]

FROM golang_deps_base as api_deps
RUN cd apps/api && go mod download
RUN cd apps/api && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as api
COPY --from=api_deps /app/apps/api/out /bin/api
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/api"]

FROM golang_deps_base as bots_deps
RUN cd apps/bots && go mod download
RUN cd apps/bots && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as bots
COPY --from=bots_deps /app/apps/bots/out /bin/bots
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/bots"]

FROM golang_deps_base as watched_deps
RUN cd apps/watched && go mod download
RUN cd apps/watched && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM go_prod_base as watched
COPY --from=watched_deps /app/apps/watched/out /bin/watched
COPY --from=base /app/docker-entrypoint.sh ./
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["/bin/watched"]
