FROM node:18-slim

RUN apt-get update && apt-get install git curl wget protobuf-compiler -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc .pnpmfile.cjs ./

COPY apps/bots apps/bots
COPY libs/prisma libs/prisma
COPY libs/shared libs/shared
COPY libs/integrations/spotify libs/integrations/spotify
COPY libs/grpc libs/grpc
COPY libs/config/ libs/config

RUN pnpm install --prod
CMD ["pnpm", "start:bots"]