FROM node:18-slim

RUN apt-get update && apt-get install git curl wget protobuf-compiler -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY node_modules node_modules
COPY apps/streamstatus apps/streamstatus
COPY libs/config/ libs/config
COPY libs/grpc libs/grpc
COPY libs/prisma libs/prisma

RUN pnpm install --prod
RUN pnpm prisma:generate

CMD ["pnpm", "start:streamstatus"]