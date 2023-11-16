FROM twir-base:builder as deps
WORKDIR /app
RUN pnpm turbo run build --filter=@twir/language-detector
RUN pnpm prune --prod

FROM twir-base:node_prod_base
WORKDIR /app
COPY --from=builder /app/apps/language-detector /app/apps/language-detector
COPY --from=builder /app/libs/grpc /app/libs/grpc
COPY --from=builder /app/node_modules /app/node_modules
CMD ["pnpm", "--filter=@twir/language-detector", "start"]
