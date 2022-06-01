FROM node:18-slim

RUN apt-get update && apt-get install git curl wget protobuf-compiler -y

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json ./
COPY apps/api/package.json apps/api/
COPY libs/prisma/package.json libs/prisma/
COPY libs/integrations/spotify/package.json libs/integrations/spotify/
COPY libs/grpc/package.json libs/grpc/
COPY libs/config/package.json libs/config/
COPY libs/shared/package.json libs/shared/

RUN pnpm install

COPY apps/api apps/api
COPY libs/prisma libs/prisma
COPY libs/integrations/spotify libs/integrations/spotify
COPY libs/grpc libs/grpc
COPY libs/config libs/config
COPY libs/shared libs/shared

RUN pnpm build:api

#FROM node:18-slim
#WORKDIR /app
#RUN npm i -g pnpm
#COPY --from=builder --chown=node:node /app/apps/backend apps/backend
#COPY --from=builder --chown=node:node /app/libs/prisma libs/prisma
#COPY package.json ./
#RUN pnpm install --prod
#RUN ls -la apps/backend/node_modules

RUN pnpm install --prod
CMD ["pnpm", "start:api"]