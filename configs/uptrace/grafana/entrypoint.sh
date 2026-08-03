#!/bin/sh
set -eu

read_secret() {
	tr -d '\r\n' < "/run/secrets/$1"
}

GF_SECURITY_ADMIN_PASSWORD="$(read_secret grafana_admin_password)"
UPTRACE_GRAFANA_TOKEN="$(read_secret uptrace_grafana_token)"
GRAFANA_POSTGRES_PASSWORD="$(read_secret grafana_postgres_password)"
GRAFANA_UPTRACE_CLICKHOUSE_PASSWORD="$(read_secret grafana_uptrace_clickhouse_password)"

export GF_SECURITY_ADMIN_PASSWORD
export UPTRACE_GRAFANA_TOKEN
export GRAFANA_POSTGRES_PASSWORD
export GRAFANA_UPTRACE_CLICKHOUSE_PASSWORD

exec /run.sh
