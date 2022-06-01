FROM node:18-slim

RUN apt-get update && apt-get install git curl wget protobuf-compiler -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./
COPY apps/streamstatus/package.json apps/streamstatus/
COPY libs/prisma/package.json libs/prisma/
COPY libs/grpc/package.json libs/grpc/
COPY libs/config/package.json libs/config/

RUN pnpm install

COPY apps/streamstatus apps/streamstatus
COPY libs/prisma libs/prisma
COPY libs/grpc libs/grpc
COPY libs/config/ libs/config

RUN pnpm build:streamstatus

RUN pnpm install --prod
CMD ["pnpm", "start:streamstatus"]