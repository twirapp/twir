#!/usr/bin/env bash

DOPPLER_TOKEN="$(cat /run/secrets/doppler_token)" doppler run -- $@
