FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/language-detector

FROM twir-base:node_prod_base as deps
WORKDIR /app
COPY --from=builder /app/apps/language-detector /app/apps/language-detector
COPY --from=builder /app/libs/grpc /app/libs/grpc
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app /app
CMD ["pnpm", "--filter=@twir/language-detector", "start"]
