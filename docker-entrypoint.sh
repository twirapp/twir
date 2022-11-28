#!/usr/bin/env sh

cat /run/secrets/doppler_token | doppler configure set token --scope /

echo "$(cat /run/secrets/doppler_token)"
echo "$@"

exec doppler run -- "$@"