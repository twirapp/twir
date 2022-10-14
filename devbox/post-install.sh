#!/bin/bash
export LC_ALL="en_US.UTF-8"
export LANG="en_US.UTF-8"

go install github.com/mitranim/gow@latest
chmod -R 700 devbox/data

pg_ctl init -D devbox/data/postgres