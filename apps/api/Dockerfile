FROM node:18-slim

RUN apt-get install git curl wget -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY apps/api apps/api
COPY libs/prisma libs/prisma
COPY libs/integrations/spotify libs/integrations/spotify
COPY libs/grpc libs/grpc
COPY libs/config libs/config
COPY libs/shared libs/shared

RUN pnpm install --prod
CMD ["pnpm", "start:api"]