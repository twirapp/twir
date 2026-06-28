FROM oven/bun:1 AS base
WORKDIR /app

FROM base AS build
COPY . .
RUN bun install --no-frozen-lockfile

FROM base AS release
WORKDIR /app
RUN apt-get update && apt-get install -y curl gnupg && \
	curl -fsSL https://cli.doppler.com/install.sh | sh -s -- --no-verify-install && \
	apt-get purge -y gnupg && apt-get autoremove -y && \
	apt-get clean && rm -rf /var/lib/apt/lists/*
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
