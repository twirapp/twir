FROM node:18-slim

RUN apt-get update && apt-get install git curl wget protobuf-compiler -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./
COPY libs/prisma/package.json libs/prisma/

RUN pnpm install

COPY libs/prisma libs/prisma

CMD ["pnpm", "migrate:deploy"]