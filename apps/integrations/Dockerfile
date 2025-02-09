FROM oven/bun:1-alpine AS base

FROM alpine:3.21 AS node_prod_base
WORKDIR /app
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/* \
COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM node_prod_base
WORKDIR /app
COPY apps/integrations/.out/twir-integrations /app
ENV NODE_ENV=production
CMD ["./twir-integrations"]
