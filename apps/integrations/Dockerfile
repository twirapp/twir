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

# bun build --compile --minify --sourcemap -e=better-sqlite3 -e=mysql -e=oracledb -e=sqlite3 -e=pg-query-stream -e=mysql2 -e=tedious src/index.js --outfile=.out/twir-integrations
FROM base
WORKDIR /app
COPY apps/integrations/.out/twir-integrations /app
ENV NODE_ENV=production
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./twir-integrations"]
