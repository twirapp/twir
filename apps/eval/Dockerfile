FROM oven/bun:1-alpine AS base

FROM alpine:3.21 AS node_prod_base
WORKDIR /app
COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM node_prod_base
WORKDIR /app
COPY apps/eval/.out/twir-eval /app
ENV NODE_ENV=production
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["./twir-eval"]
