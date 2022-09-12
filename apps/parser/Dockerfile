FROM golang:1.19-alpine AS build_base

RUN apk add --no-cache git

WORKDIR /tmp/go-sample-app

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum /apps/parser/

RUN cd /apps/parser && go mod download

COPY /apps/parser /apps/parser
COPY /libs/nats libs/nats

RUN go build -o ./out/parser /apps/parser/cmd/main.go

# Start fresh from a smaller image
FROM alpine:3.9 
RUN apk add ca-certificates

COPY --from=build_base /tmp/go-sample-app/out/parser /app/parser

CMD ["/app/parser"]