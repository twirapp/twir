#!/usr/bin/env sh


cat /run/secrets/doppler_token | doppler configure set token --scope /

exec doppler run -- "$@"