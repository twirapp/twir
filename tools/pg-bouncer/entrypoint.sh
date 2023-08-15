#!/usr/bin/env sh

export POSTGRESQL_USERNAME=$(cat /run/secrets/tsuwari_postgres_user)
export POSTGRESQL_PASSWORD=$(cat /run/secrets/tsuwari_postgres_password)
export POSTGRESQL_DATABASE=$(cat /run/secrets/tsuwari_postgres_db)
export POSTGRESQL_HOST=postgres
