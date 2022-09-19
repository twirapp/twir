#!/bin/bash
cd /app/apps/parser && go mod download
cd ../..
pnpm install && pnpm build:backend && pnpm dev