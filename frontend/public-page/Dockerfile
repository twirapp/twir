FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/public-page

FROM devforth/spa-to-http:latest
COPY --from=builder /app/frontend/public-page/dist/ .
