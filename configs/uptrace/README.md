# Uptrace on the Twir Docker Swarm

The stack runs the Uptrace UI/API on application nodes and keeps its dedicated ClickHouse,
PostgreSQL, and Redis services on the node labeled `databases=true`. `node-exporter` and cAdvisor
run globally; one OpenTelemetry Collector discovers and scrapes every exporter task.

## One-time preparation

Ensure the database node has the label used by the stack:

```sh
docker node update --label-add databases=true DATABASE_NODE
```

Create the external Docker secrets. Secret values must be single-line strings; hex is recommended
for generated infrastructure secrets. Preserve the project token from the existing Uptrace server
so existing SDK configuration keeps working:

```sh
project_token="$(ssh root@169.58.78.216 "sed -n 's/^PROJECT_TOKEN=//p' /opt/uptrace/.credentials")"
printf 'http://%s@uptrace:4317/1' "$project_token" | docker secret create uptrace_dsn -
unset project_token

openssl rand -hex 32 | docker secret create uptrace_service_secret -
openssl rand -hex 32 | docker secret create uptrace_clickhouse_password -
openssl rand -hex 32 | docker secret create uptrace_postgres_password -

read -r -s -p 'Uptrace admin password: ' uptrace_admin_password
printf '%s' "$uptrace_admin_password" | docker secret create uptrace_admin_password -
unset uptrace_admin_password
```

The internal DSN above is intentionally used by the Collector. Applications inside the Swarm can
send unauthenticated OTLP to `otel-collector:4317`; the Collector adds the Uptrace DSN when it
forwards telemetry. Set the application configuration to:

```text
OTEL_ENDPOINT=otel-collector:4317
OTEL_INSECURE=true
OTEL_HEADERS=
```

## Deploy

```sh
docker stack deploy --with-registry-auth -c docker-compose.stack.yml twir
```

Uptrace is routed through Traefik at `https://uptrace.twir.app`. The separate Grafana service is no
longer needed: Uptrace provides the UI and is Prometheus/Tempo compatible if Grafana is added again
later. If an old Grafana service still belongs to the `twir` stack, verify its name and remove it
after Uptrace is healthy, or use `docker stack deploy --prune` during the planned cutover.

## Verify

```sh
docker stack services twir
docker service logs twir_uptrace --since 10m
docker service logs twir_otel-collector --since 10m
docker service ps twir_node-exporter
docker service ps twir_cadvisor
```

The `uptrace-clickhouse-ttl` one-shot service waits for Uptrace migrations and then applies the same
seven-day telemetry retention policy used by the standalone server. The ClickHouse config also
disables high-volume system log tables and retains `system.query_log` for one day.
