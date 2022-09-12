FROM ubuntu:latest AS build_base

RUN apt-get update && apt-get install -y curl git protobuf-compiler wget
RUN curl -fsSL https://deb.nodesource.com/setup_18.x | bash -
RUN apt-get install nodejs

RUN wget https://go.dev/dl/go1.19.1.linux-amd64.tar.gz && rm -rf /usr/local/go && tar -C /usr/local -xzf go1.19.1.linux-amd64.tar.gz && export PATH=$PATH:/usr/local/go/bin

ENV PATH="/usr/local/go/bin:${PATH}" 

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

RUN npm i -g pnpm@7
WORKDIR /app

COPY apps/parser/go.mod apps/parser/go.sum /apps/parser/
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY /apps/parser /apps/parser/
COPY /libs/nats /libs/nats/

RUN pnpm --filter=@tsuwari/nats build
RUN cd /apps/parser && go mod download

RUN cd /apps/parser && go build -o /app/out/parser ./cmd/main.go
RUN ls /apps/parser

FROM alpine:3.16 
RUN apk add ca-certificates

WORKDIR /app

COPY --from=build_base /app/out/parser /app/parser

CMD ["/app/parser"]