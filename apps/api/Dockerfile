FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/language-detector
RUN cd apps/api && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /twir_application ./cmd/main.go && upx -9 -k /twir_application

FROM twir-base:go_prod_base
COPY --from=builder /app/twir_application /bin/twir_application
CMD ["/bin/twir_application"]
