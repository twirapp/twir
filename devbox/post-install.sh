#!/bin/bash
export LC_ALL="en_US.UTF-8"
export LANG="en_US.UTF-8"

go install github.com/mitranim/gow@latest
go install mvdan.cc/gofumpt@latest
go install github.com/segmentio/golines@latest

echo "Adminer available at http://localhost:8080/?pgsql=postgres:5432&username=tsuwari"