FROM mirror.gcr.io/alpine:latest as go_prod_base
WORKDIR /app
RUN apk add wget gzip && \
    wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add doppler && \
    mkdir -p assets/geo-db && \
    YEAR=$(date +%Y) && \
    MONTH=$(date +%m) && \
    wget -q -t3 -O /tmp/geo-db.mmdb.gz "https://download.db-ip.com/free/dbip-city-lite-${YEAR}-${MONTH}.mmdb.gz" && \
    gunzip -c /tmp/geo-db.mmdb.gz > assets/geo-db/dbip-city-lite.mmdb && \
    rm -f /tmp/geo-db.mmdb.gz && \
    echo "{\"lastUpdated\":\"$(date -Iseconds)\",\"version\":\"build-$(date +%Y%m%d)\",\"recordsCount\":null}" > assets/geo-db/db-info.json && \
    apk del wget gzip && \
    rm -rf /var/cache/apk/*
COPY docker-entrypoint.sh /app/
RUN chmod +x /app/docker-entrypoint.sh
ENTRYPOINT ["/app/docker-entrypoint.sh"]

FROM go_prod_base
WORKDIR /app
COPY --from=go_prod_base /app/assets /app/assets
COPY apps/api-gql/.out/twir-api-gql /app/twir-api-gql
CMD ["/app/twir-api-gql"]
