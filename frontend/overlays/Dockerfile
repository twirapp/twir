FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/overlays

FROM caddy:latest
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /app/frontend/overlays/dist/ /app
EXPOSE 80
