FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/landing
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=builder /app/frontend/landing /app/frontend/landing
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/grpc /app/libs/grpc
COPY --from=builder /app/node_modules /app/node_modules
CMD ["pnpm", "--filter=@twir/landing", "start"]
