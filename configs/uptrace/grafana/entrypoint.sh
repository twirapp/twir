#!/bin/sh
set -eu

read_secret() {
	tr -d '\r\n' < "/run/secrets/$1"
}

uptrace_dsn="$(read_secret uptrace_dsn)"
case "$uptrace_dsn" in
	*://*@*/[0-9]*) ;;
	*)
		echo "uptrace_dsn must look like http://TOKEN@uptrace:4317/1" >&2
		exit 1
		;;
esac

UPTRACE_PROJECT_TOKEN="${uptrace_dsn#*://}"
UPTRACE_PROJECT_TOKEN="${UPTRACE_PROJECT_TOKEN%%@*}"
GF_SECURITY_ADMIN_PASSWORD="$(read_secret grafana_admin_password)"

export UPTRACE_PROJECT_TOKEN
export GF_SECURITY_ADMIN_PASSWORD

exec /run.sh
