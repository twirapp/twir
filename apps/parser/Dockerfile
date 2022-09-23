FROM ubuntu:latest AS build_base

RUN apt-get update && apt-get install -y curl git protobuf-compiler wget
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install nodejs

RUN wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz

ENV GOPATH=$HOME/go
ENV PATH=$PATH:$GOPATH/bin:/usr/local/go/bin

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

RUN npm i -g pnpm@7
WORKDIR /app

COPY apps/parser/go.mod apps/parser/go.sum ./apps/parser/
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./
COPY /libs/nats/package.json ./libs/nats/
RUN pnpm install

COPY /apps/parser ./apps/parser/
COPY /libs/nats ./libs/nats/
COPY /libs/config ./libs/config/
COPY /libs/gomodels ./libs/gomodels/

RUN pnpm --filter=@tsuwari/nats build
RUN cd ./apps/parser && go mod download

RUN cd ./apps/parser && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out ./cmd/main.go

FROM alpine:latest
RUN apk add --no-cache ca-certificates
RUN wget -q -t3 'https://packages.doppler.com/public/cli/rsa.8004D9FF50437357.key' -O /etc/apk/keys/cli@doppler-8004D9FF50437357.rsa.pub && \
    echo 'https://packages.doppler.com/public/cli/alpine/any-version/main' | tee -a /etc/apk/repositories && \
    apk add --no-cache doppler

COPY --from=build_base /app/apps/parser/out /bin/parser
RUN chmod +x /bin/parser

RUN rm -rf /var/cache/apk/*

ENTRYPOINT ["doppler", "run", "--"]
CMD ["/bin/parser"]