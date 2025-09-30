FROM mirror.gcr.io/oven/bun:1-alpine AS base
WORKDIR /app

FROM base AS builder
COPY . .
ENV NODE_ENV=production
RUN bun install --frozen-lockfile
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
CMD ["bun", "--bun", "run", "server/index.mjs"]
