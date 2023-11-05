FROM twir-base:builder as builder

RUN cd apps/timers && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /twir_application ./cmd/main.go && upx -9 -k /twir_application

FROM twir-base:go_prod_base
COPY --from=builder /app/twir_application /bin/twir_application
CMD ["/bin/twir_application"]
