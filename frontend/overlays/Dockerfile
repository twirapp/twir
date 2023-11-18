FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/overlays

FROM devforth/spa-to-http:latest
COPY --from=builder /app/frontend/overlays/dist/ .
