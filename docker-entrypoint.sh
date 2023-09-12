#!/usr/bin/env sh

cat /run/secrets/twir_doppler_token | doppler configure set token --scope / > /dev/null

exec doppler run -- "$@"
