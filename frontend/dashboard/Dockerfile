FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/dashboard

FROM caddy:latest
COPY Caddyfile /etc/caddy/Caddyfile
COPY --from=builder /app/frontend/dashboard/dist/ /app
EXPOSE 80
