FROM twir-base:builder as deps
WORKDIR /app
RUN pnpm turbo run build --filter=@twir/eval
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app/apps/eval /app/apps/eval
COPY --from=deps /app/libs/config /app/libs/config
COPY --from=deps /app/libs/grpc /app/libs/grpc
COPY --from=deps /app/node_modules /app/node_modules
ENTRYPOINT ["/bin/sh", "-c"]
CMD ["pnpm --filter=@twir/eval start"]
