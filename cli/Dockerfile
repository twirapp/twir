FROM mirror.gcr.io/golang:1.26-alpine AS builder
WORKDIR /src

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/twir-cli ./cli/main.go

FROM mirror.gcr.io/alpine:3.22
WORKDIR /app
COPY --from=builder /out/twir-cli /app/twir-cli

ENTRYPOINT ["/app/twir-cli"]
CMD ["deploy-webhook"]
