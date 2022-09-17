FROM node:18-slim

ENV LAST_UPDATED 2022-06-19

RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg && \
  curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | apt-key add - && \
  echo "deb https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list

RUN apt-get update && apt-get install -y git curl wget doppler

WORKDIR /app

RUN npm i -g pnpm

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./
COPY libs/typeorm/package.json libs/typeorm/

COPY libs/config libs/config

RUN pnpm install
COPY libs/typeorm libs/typeorm

CMD ["pnpm", "run", "migrate:deploy"]