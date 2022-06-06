FROM node:18-slim

RUN apt-get update && apt-get install git curl wget -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY node_modules node_modules
COPY apps/scheduler apps/scheduler
COPY libs/prisma libs/prisma
COPY libs/grpc libs/grpc
COPY libs/config/ libs/config

RUN pnpm install --prod
RUN pnpm prisma:generate

CMD ["pnpm", "start:scheduler"]