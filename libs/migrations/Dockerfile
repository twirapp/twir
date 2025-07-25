FROM mirror.gcr.io/alpine:latest as go_prod_base
WORKDIR /app
RUN apk add wget && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && apk del wget && \
    rm -rf /var/cache/apk/*
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM mirror.gcr.io/golang:1.24.1-alpine as builder
WORKDIR /app
COPY . .
RUN cd cli && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/twir_application ./main.go

FROM go_prod_base
COPY --from=builder /app/twir_application /app/twir_application
COPY libs/migrations /app/libs/migrations
CMD ["sh", "-c", "/app/twir_application m run && sleep 60"]
