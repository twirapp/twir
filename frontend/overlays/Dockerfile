FROM node:18-alpine as builder

WORKDIR /app

RUN apk add git openssh libc6-compat
RUN npm i -g pnpm@7

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml tsconfig.base.json tsconfig.json turbo.json .npmrc ./

COPY frontend frontend
COPY patches patches

RUN pnpm install --frozen-lockfile
RUN cd frontend/overlays && pnpm build

FROM codecentric/single-page-application-server
COPY --from=builder /app/frontend/overlays/dist/ /app