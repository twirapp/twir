FROM oven/bun:1-alpine AS base
WORKDIR /app

FROM base AS installer
COPY package.json bun.lock bunfig.toml .npmrc ./
COPY libs/brand/package.json ./libs/brand/package.json
COPY libs/config/package.json ./libs/config/package.json
COPY web/package.json ./web/package.json
RUN bun install --no-save

FROM base AS builder
COPY --from=installer /app/node_modules node_modules
COPY package.json bun.lock bunfig.toml .npmrc ./
COPY libs/brand ./libs/brand
COPY libs/config ./libs/config
COPY web ./web
COPY apps/api-gql/schema ./apps/api-gql/schema

ENV NODE_ENV=production
RUN bun --filter=@twir/web run build

FROM base AS prod-base
ENV NODE_ENV=production
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler
COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM prod-base
COPY --from=builder /app/web/.output ./
CMD ["bun", "run", "server/index.mjs"]
