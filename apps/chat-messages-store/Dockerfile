FROM node:20-alpine as builder
WORKDIR /app
RUN npm i -g pnpm@8

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY apps/chat-messages-store /app/apps/chat-messages-store
COPY libs/config /app/libs/config
COPY libs/bus-core /app/libs/bus-core

RUN pnpm prune --prod

FROM node:20-alpine as node_prod_base
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

FROM node_prod_base
WORKDIR /app
COPY --from=builder /app /app
CMD ["pnpm", "--filter=@twir/chat-messages-store", "start"]
