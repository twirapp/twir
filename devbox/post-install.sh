#!/bin/bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

which gow || go install github.com/mitranim/gow@latest
which gofumpt || go install mvdan.cc/gofumpt@latest
which golines || go install github.com/segmentio/golines@latest
which twitch-cli || go install github.com/twitchdev/twitch-cli@latest

which nodemon || pnpm add -g nodemon

echo "Adminer available at http://localhost:8080/?pgsql=postgres:5432&username=tsuwari"