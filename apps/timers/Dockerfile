FROM golang:1.20.1-alpine as builder

WORKDIR /app

RUN apk add git curl wget upx

COPY go.work go.work.sum docker-entrypoint.sh ./
COPY libs libs
COPY apps/timers apps/timers

RUN cd apps/timers && go mod download
RUN cd apps/timers && CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o ./out ./cmd/main.go && upx -9 -k ./out

FROM alpine:latest
RUN apk add wget && \
  wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
  echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
  apk add doppler && apk del wget && \
  rm -rf /var/cache/apk/*
COPY --from=builder /app/apps/timers/out /bin/timers
COPY --from=builder /app/docker-entrypoint.sh /app/
ENTRYPOINT ["/app/docker-entrypoint.sh"]
CMD ["/bin/timers"]