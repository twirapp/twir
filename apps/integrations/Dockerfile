FROM twir-base:builder as builder
WORKDIR /app
RUN pnpm turbo run build --filter=@twir/grpc
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=builder /app/apps/integrations /app/apps/integrations
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/grpc /app/libs/grpc
COPY --from=builder /app/libs/pubsub /app/libs/pubsub
COPY --from=builder /app/node_modules /app/node_modules
CMD ["pnpm", "--filter=@twir/integrations", "start"]
