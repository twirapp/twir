#!/usr/bin/env sh

cat /run/secrets/tsuwari_doppler_token | doppler configure set token --scope / > /dev/null

exec doppler run -- "$@"