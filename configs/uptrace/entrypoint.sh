#!/bin/sh
set -eu

read_secret() {
	tr -d '\r\n' < "/run/secrets/$1"
}

escape_yaml_sed() {
	printf '%s' "$1" | sed -e 's/\\/\\\\/g' -e 's/"/\\"/g' -e 's/[\\|&]/\\&/g'
}

uptrace_dsn="$(read_secret uptrace_dsn)"
case "$uptrace_dsn" in
	*://*@*/[0-9]*) ;;
	*)
		echo "uptrace_dsn must look like http://TOKEN@uptrace:4317/1" >&2
		exit 1
		;;
esac

project_token="${uptrace_dsn#*://}"
project_token="${project_token%%@*}"
service_secret="$(read_secret uptrace_service_secret)"
admin_password="$(read_secret uptrace_admin_password)"
clickhouse_password="$(read_secret uptrace_clickhouse_password)"
postgres_password="$(read_secret uptrace_postgres_password)"

sed \
	-e "s|__UPTRACE_PROJECT_TOKEN__|$(escape_yaml_sed "$project_token")|g" \
	-e "s|__UPTRACE_SERVICE_SECRET__|$(escape_yaml_sed "$service_secret")|g" \
	-e "s|__UPTRACE_ADMIN_PASSWORD__|$(escape_yaml_sed "$admin_password")|g" \
	-e "s|__UPTRACE_CLICKHOUSE_PASSWORD__|$(escape_yaml_sed "$clickhouse_password")|g" \
	-e "s|__UPTRACE_POSTGRES_PASSWORD__|$(escape_yaml_sed "$postgres_password")|g" \
	/etc/uptrace/config.yml.tmpl > /tmp/uptrace.yml
chmod 600 /tmp/uptrace.yml

retry() {
	operation="$1"
	shift
	attempt=1
	while ! "$@"; do
		if [ "$attempt" -ge 60 ]; then
			echo "$operation failed after $attempt attempts" >&2
			exit 1
		fi
		attempt=$((attempt + 1))
		sleep 5
	done
}

retry "PostgreSQL initialization" /uptrace --config=/tmp/uptrace.yml pg init
/uptrace --config=/tmp/uptrace.yml pg migrate
retry "ClickHouse initialization" /uptrace --config=/tmp/uptrace.yml ch init
/uptrace --config=/tmp/uptrace.yml ch migrate
/uptrace --config=/tmp/uptrace.yml db seed
exec /uptrace --config=/tmp/uptrace.yml serve
