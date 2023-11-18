FROM twir-base:builder as builder
RUN pnpm turbo run build --filter=@twir/grpc
RUN cd apps/parser && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/twir_application ./cmd/main.go && \
    upx -9 -k /app/twir_application

FROM twir-base:go_prod_base
COPY --from=builder /app/twir_application /bin/twir_application
CMD ["/bin/twir_application"]
