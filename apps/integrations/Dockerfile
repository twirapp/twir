FROM twir-base:builder as deps
WORKDIR /app
RUN pnpm turbo run build --filter=@twir/integrations
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app/apps/integrations /app/apps/integrations
COPY --from=deps /app/libs/config /app/libs/config
COPY --from=deps /app/libs/grpc /app/libs/grpc
COPY --from=deps /app/libs/pubsub /app/libs/pubsub
COPY --from=deps /app/node_modules /app/node_modules
CMD ["pnpm", "--filter=@twir/integrations", "start"]
