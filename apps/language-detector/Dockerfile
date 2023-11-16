FROM twir-base:builder as deps
WORKDIR /app
COPY --from=builder /app/apps/language-detector /app/apps/language-detector
COPY --from=builder /app/libs/grpc /app/libs/grpc
RUN pnpm turbo run build --filter=@twir/language-detector
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=deps /app /app
CMD ["pnpm", "--filter=@twir/language-detector", "start"]
