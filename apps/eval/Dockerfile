FROM oven/bun:1 AS base

FROM base AS node_prod_base
WORKDIR /app
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/*
COPY package.json .npmrc bunfig.toml docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM base
WORKDIR /app
COPY apps/eval/.out/twir-eval /app
ENV NODE_ENV=production
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./twir-eval"]
