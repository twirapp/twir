FROM node:18-alpine as builder

COPY --from=golang:1.20.1-alpine /usr/local/go/ /usr/local/go/
ENV PATH="$PATH:/usr/local/go/bin"
ENV PATH="$PATH:/root/go/bin"

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1 && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app

RUN apk add git openssh libc6-compat protoc
RUN npm i -g pnpm@7

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY libs libs
COPY apps/streamstatus apps/streamstatus
COPY patches patches

RUN pnpm install --frozen-lockfile
RUN pnpm build:libs
RUN cd apps/streamstatus && pnpm build

FROM node:18-alpine
RUN apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*

RUN npm i -g pnpm@7

WORKDIR /app

COPY --from=builder /app/package.json /app/pnpm-lock.yaml /app/pnpm-workspace.yaml /app/turbo.json /app/.npmrc ./
COPY --from=builder /app/apps/streamstatus apps/streamstatus/
COPY --from=builder /app/libs/config libs/config/
COPY --from=builder /app/libs/shared libs/shared/
COPY --from=builder /app/libs/typeorm libs/typeorm/
COPY --from=builder /app/libs/pubsub libs/pubsub/
COPY --from=builder /app/patches patches
RUN pnpm prune --prod

ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["pnpm", "--filter=@tsuwari/streamstatus", "start"]