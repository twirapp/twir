FROM oven/bun:1-alpine AS base
WORKDIR /app

FROM base AS build
COPY . .
RUN bun install --frozen-lockfile

FROM base AS release
WORKDIR /app
RUN apk add wget libstdc++ libc6-compat && \
	wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
	echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
	apk add doppler && apk del wget && \
	rm -rf /var/cache/apk/*
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
COPY --from=build /app/node_modules /app/node_modules
COPY --from=build /app/apps/executron /app/apps/executron
COPY --from=build /app/libs /app/libs
COPY --from=build /app/package.json /app/package.json
WORKDIR /app/apps/executron
ENV NODE_ENV=production
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["bun", "run", "src/index.ts"]
