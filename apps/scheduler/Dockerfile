FROM node:18-slim

ENV LAST_UPDATED 2022-06-19

RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg && \
  curl -sLf --retry 3 --tlsv1.2 --proto "=https" 'https://packages.doppler.com/public/cli/gpg.DE2A7741A397C129.key' | apt-key add - && \
  echo "deb https://packages.doppler.com/public/cli/deb/debian any-version main" | tee /etc/apt/sources.list.d/doppler-cli.list

RUN apt-get update && apt-get install git curl wget -y doppler

WORKDIR /app

RUN npm i -g pnpm@7.3.0

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY node_modules node_modules
COPY apps/scheduler apps/scheduler
COPY libs/prisma libs/prisma
COPY libs/config/ libs/config
COPY libs/shared/ libs/shared

RUN pnpm install --prod
RUN pnpm prisma:generate

ENTRYPOINT ["doppler", "run", "--"]
CMD ["pnpm", "start:scheduler"]