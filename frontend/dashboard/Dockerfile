FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/dashboard

FROM devforth/spa-to-http:latest
COPY --from=builder /app/frontend/dashboard/dist/ .
