FROM node:22-alpine AS base
WORKDIR /app

FROM base AS builder

RUN npm install -g turbo@1

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml .npmrc ./
COPY libs libs
COPY web web

RUN turbo prune --scope=@twir/web --docker

FROM base as installer

RUN npm install -g pnpm@9
RUN npm install -g turbo@1

COPY --from=builder /app/out/json/ .
COPY --from=builder /app/out/pnpm-lock.yaml ./pnpm-lock.yaml

RUN pnpm install

COPY --from=builder /app/out/full/ .
COPY apps/api-gql/schema ./apps/api-gql/schema
COPY tsconfig.json tsconfig.json
COPY turbo.json turbo.json

RUN NODE_ENV=production turbo run build --filter=@twir/web
RUN turbo prune --scope=@twir/web --docker
RUN pnpm prune --prod

FROM base as prod-base
ENV NODE_ENV production
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/* \
COPY docker-entrypoint.sh ./
RUN chmod +x docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM base AS prod-base
COPY --from=installer /app/out/full .
CMD ["node", "web/.output/server/index.mjs"]
