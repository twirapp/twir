FROM twir-base:builder as builder
RUN cd libs/migrations && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./twir_application ./main.go && upx -9 -k ./twir_application

FROM twir-base:go_prod_base
COPY --from=builder /app/twir_application /bin/twir_application
CMD ["sh", "-c", "/bin/twir_application && sleep 60"]
