FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/public-page

FROM caddy:latest
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /app/frontend/public-page/dist/ /app
EXPOSE 80
