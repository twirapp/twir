# Uptrace and Grafana on the Twir Docker Swarm

The stack runs the Uptrace UI/API on application nodes and keeps its dedicated ClickHouse,
PostgreSQL, and Redis services on the node labeled `databases=true`. `node-exporter` and cAdvisor
run globally; one OpenTelemetry Collector discovers and scrapes every exporter task.

## One-time preparation

Ensure the database node has the label used by the stack:

```sh
docker node update --label-add databases=true DATABASE_NODE
```

Create the external Docker secrets on a Swarm manager. Secret values must be single-line strings;
hex is recommended for generated infrastructure secrets. Applications send telemetry through the
Collector, so the Uptrace installation uses a project DSN:

```sh
project_token="$(openssl rand -hex 32)"
printf 'http://%s@uptrace:4317/1' "$project_token" | docker secret create uptrace_dsn -
unset project_token

openssl rand -hex 32 | docker secret create uptrace_service_secret -
openssl rand -hex 32 | docker secret create uptrace_clickhouse_password -
openssl rand -hex 32 | docker secret create uptrace_postgres_password -

uptrace_admin_password="$(openssl rand -hex 16)"
printf 'Save Uptrace admin password: %s\n' "$uptrace_admin_password"
printf '%s' "$uptrace_admin_password" | docker secret create uptrace_admin_password -
unset uptrace_admin_password

grafana_admin_password="$(openssl rand -hex 16)"
printf 'Save Grafana admin password: %s\n' "$grafana_admin_password"
printf '%s' "$grafana_admin_password" | docker secret create grafana_admin_password -
unset grafana_admin_password
```

Grafana needs three additional secrets:

- `uptrace_grafana_token`: an Uptrace user token, not the project token from `uptrace_dsn`;
- `grafana_postgres_password`: password of the read-only `grafana_ro` PostgreSQL role;
- `grafana_uptrace_clickhouse_password_v2`: password of the read-only `grafana_ro` ClickHouse
  user.

Create them without putting values in shell history:

```sh
printf 'Uptrace user token: '
read -r -s value
printf '\n'
printf '%s' "$value" | docker secret create uptrace_grafana_token -
unset value

value="$(openssl rand -hex 32)"
printf '%s' "$value" | docker secret create grafana_postgres_password -
printf 'Store this value securely for the PostgreSQL role: %s\n' "$value"
unset value

value="$(openssl rand -hex 32)"
printf '%s' "$value" | docker secret create grafana_uptrace_clickhouse_password_v2 -
printf 'Store this value securely for the ClickHouse user: %s\n' "$value"
unset value
```

Create `grafana_ro` in the main PostgreSQL database with `LOGIN`, `CONNECT` to `twir`, `USAGE` on
`public`, and `SELECT` only on `users`, `channels`, `channel_platforms`, `channels_streams`, and
`users_stats`. The last table powers the same `Messages processed` total as the Twir landing page.
Set these role defaults:

```sql
GRANT SELECT ON TABLE users, channels, channel_platforms, channels_streams, users_stats TO grafana_ro;
ALTER ROLE grafana_ro SET default_transaction_read_only = on;
ALTER ROLE grafana_ro SET statement_timeout = '15s';
```

Create `grafana_ro` in Uptrace ClickHouse, grant only `SELECT ON uptrace.spans_index`, and assign a
settings profile with `readonly = 1`, `max_execution_time = 30`, and
`max_memory_usage = 536870912`. `CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1` in the stack grants the
Uptrace administrative user the access-management privileges needed for this one-time setup.

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

Uptrace is routed through Traefik at `https://uptrace.twir.app`. Grafana is available at
`https://grafana.twir.app` and is provisioned with four data sources:

- Uptrace Prometheus for node-exporter and cAdvisor metrics;
- Uptrace Tempo for traces;
- Uptrace ClickHouse for EventSub and HTTP telemetry stored in `spans_index`;
- Twir PostgreSQL for aggregate user, channel, and live-stream statistics.

The provisioned `Twir Live Overview` dashboard is stored at
`configs/uptrace/grafana/dashboards/twir-overview.json`. Its application panels use OTel spans from
Uptrace for command/handler p95 and PostgreSQL, ClickHouse, and Redis read/write rates. These rates
describe operations observed from Twir applications, not every server-side background query. CPU
and RAM are grouped by Swarm node from node-exporter.

Three more provisioned dashboards cover infrastructure metrics:

- `Twir Nodes` (`twir-nodes.json`) — per-node CPU, load, memory, root filesystem, disk I/O and
  IOPS, network throughput, file descriptors, and conntrack usage from node-exporter, with a
  `Node` variable. Series are joined to `node_uname_info` for readable node names.
- `Twir Containers` (`twir-containers.json`) — per-container CPU (including CFS throttling),
  memory usage and % of limit, network, and disk I/O from cAdvisor, filterable by `Service` and
  `Container` variables.
- `Twir Swarm Services` (`twir-swarm.json`) — cluster stats aggregated per Swarm service
  (replicas, CPU, memory, % of limit, network, disk I/O). The `Service` variable selects one or
  more services; the bottom section then drills into individual tasks of the selection
  (CPU/memory/network/throttling per task plus task age for spotting crash loops).

The container and swarm dashboards rely on cAdvisor's default label export
(`--store_container_labels=true`, the default since long before v0.57): Swarm automatically stamps
task containers with `com.docker.swarm.service.name` / `com.docker.swarm.task.name`, which appear
in Prometheus as `container_label_com_docker_swarm_service_name` and
`container_label_com_docker_swarm_task_name`. No extra cAdvisor flags are needed.

New dashboards are added as new Swarm config keys (e.g. `uptrace_grafana_dashboard_nodes`) mounted
into `/etc/grafana/dashboards/`; the file provider picks them up on service update. As with any
other config, once a dashboard config exists, edits require a new key (see the note below).

The streamer opens the dashboard over the attachable
`twir` overlay network using `http://twir_grafana:3000`, so it does not depend on Cloudflare or
public DNS. The stack-qualified service name avoids the legacy `grafana` alias from another stack.
Grafana is pinned to `traefik-1` because `grafana-data` is a node-local volume; moving the task to
another node would create a separate SQLite database and a different set of users/sessions.

Docker Swarm configs are immutable. When changing an existing datasource or entrypoint config,
create a new key such as `uptrace_grafana_datasource_v3`, point the service at it, and deploy. Do
not try to overwrite the old config name.

Apply the committed PostgreSQL migration before relying on the 30-day new-users panel:

```sh
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_users_created_at ON users (created_at);
```

The `Users` stat runs an exact `count(*)` on the users table (the previous
`pg_class.reltuples` estimate only moved when ANALYZE ran and looked frozen). On multi-million-row
tables this takes on the order of a second; the `grafana_ro` role already has `SELECT` on `users`
and a 15s statement timeout, which is sufficient.

For a fresh streamer host:

```sh
cd /root/services/twir-streamer
git pull --ff-only origin main
# Set grafanaEnabled=true, the internal dashboard URL, and a dedicated Viewer account in config.json.
docker compose up -d --build
```

## Verify

```sh
docker stack services twir
docker service logs twir_uptrace --since 10m
docker service logs twir_otel-collector --since 10m
docker service logs twir_grafana --since 10m
docker service ps twir_node-exporter
docker service ps twir_cadvisor
```

Also verify datasource health and the dashboard with the Grafana API or UI, then check connectivity
from the streamer container:

```sh
docker compose exec streamer bun -e \
  'fetch("http://twir_grafana:3000/api/health").then(async r => { console.log(r.status, await r.text()) })'
```

The `uptrace-clickhouse-ttl` one-shot service waits for Uptrace migrations and then applies the same
seven-day telemetry retention policy used by the standalone server. The ClickHouse config also
disables high-volume system log tables and retains `system.query_log` for one day.
