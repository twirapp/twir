FROM twir-base:node_prod_base as deps
WORKDIR /app
COPY --from=builder /app/apps/eval /app/apps/eval
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/grpc /app/libs/grpc
RUN pnpm turbo run build --filter=@twir/eval
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app /app
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["pnpm --filter=@twir/eval start"]
