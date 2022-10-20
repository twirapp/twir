#!/bin/bash
go install github.com/mitranim/gow@latest
go install mvdan.cc/gofumpt@latest
go install github.com/segmentio/golines@latest

echo "Adminer available at http://localhost:8080/?pgsql=postgres:5432&username=tsuwari"