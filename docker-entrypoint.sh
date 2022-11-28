#!/usr/bin/env sh

exec DOPPLER_TOKEN="$(cat /run/secrets/doppler_token)" doppler run -- "$@"