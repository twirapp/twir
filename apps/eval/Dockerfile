FROM twir-base:builder as builder
WORKDIR /app
RUN pnpm turbo run build --filter=@twir/eval
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=builder /app/apps/eval /app/apps/eval
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/grpc /app/libs/grpc
COPY --from=builder /app/node_modules /app/node_modules
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["pnpm --filter=@twir/eval start"]
