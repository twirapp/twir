FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/integrations

FROM twir-base:node_prod_base as deps
WORKDIR /app
COPY --from=builder /app/apps/integrations /app/apps/integrations
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/grpc /app/libs/grpc
COPY --from=builder /app/libs/pubsub /app/libs/pubsub
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app /app
CMD ["pnpm", "--filter=@twir/integrations", "start"]
