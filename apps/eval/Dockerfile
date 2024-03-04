FROM twir-base:builder as builder
WORKDIR /app
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=builder /app/apps/eval /app/apps/eval
COPY --from=builder /app/libs/config /app/libs/config
COPY --from=builder /app/libs/bus-core /app/libs/bus-core
COPY --from=builder /app/node_modules /app/node_modules
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["pnpm --filter=@twir/eval start"]
