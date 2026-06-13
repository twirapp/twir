# Twir Kubernetes Manifests

This directory contains deployable Kubernetes manifests for the Twir stack represented by `docker-compose.stack.yml`.

The production entry point is:

```bash
kubectl kustomize deploy/kubernetes/overlays/prod
kubectl apply -k deploy/kubernetes/overlays/prod
```

`kubectl` and standalone `kustomize` both work with this layout. The manifests target k3s with the built-in Traefik ingress controller and `local-path` storage.

## Important Production Notes

- Namespace is always `twir`.
- These manifests do not provide database HA. Postgres, ClickHouse, Redis, Temporal Postgres, and NATS JetStream use single-replica StatefulSets with local PVCs.
- Most Twir images are still `latest` to match compose. Replace them with immutable tags or digests before production cutover.
- No real secrets are committed. Runtime secrets must be created from your own env file.
- HAProxy from compose is represented by Traefik Ingress and Traefik middlewares. Standard Ingress cannot fully reproduce HAProxy stick-table rate limiting, cache, gzip, and custom-domain fallback rewrite. Keep HAProxy as an edge DaemonSet outside this kustomize tree if exact fallback/rate-limit behavior is required.
- The included OTel config is Kubernetes-safe and does not mount the Docker socket. Replace the debug exporter with your real OTLP backend before production.

## Layout

- `base/namespace.yaml`: `twir` namespace.
- `base/service-accounts.yaml`: workload service accounts and image pull secret reference.
- `base/configmaps.yaml`: ClickHouse XML and PgDog config/start script.
- `base/secrets.example.env`: keys required for the runtime `twir-secrets` Secret.
- `base/secret-template.yaml`: intentionally unreferenced Secret skeleton; do not fill it in git.
- `base/storage.yaml`: PVCs using k3s `local-path`.
- `base/stateful.yaml`: Postgres, ClickHouse, Redis, Temporal Postgres, NATS.
- `base/poolers.yaml`: PgDog and PgBouncer.
- `base/apps.yaml`: Twir app workloads, Adminer, Temporal, Temporal UI, TTS.
- `base/jobs.yaml`: migrations Job and Postgres backup CronJob.
- `base/ingress.yaml`: k3s Traefik ingress replacement for public routes.
- `base/observability.yaml`: OTel collector Deployment and config.
- `overlays/prod/kustomization.yaml`: production render target.

## Node Labels

The manifests use the labels below to match the Swarm placement intent.

```bash
kubectl label node twir-k3s-1 twir.app/database-primary=true twir.app/database=true --overwrite
kubectl label node twir-k3s-2 twir.app/database=true twir.app/cache=true twir.app/backup=true --overwrite
kubectl label node twir-k3s-3 twir.app/apps=true --overwrite
```

Placement summary:

- Postgres, ClickHouse, and Temporal Postgres require `twir.app/database-primary=true`.
- Redis requires `twir.app/cache=true`.
- Postgres backup requires `twir.app/backup=true`.
- Language processor and toxicity detector require `twir.app/apps=true`.
- Other app workloads prefer nodes without `twir.app/database` but can schedule elsewhere if capacity requires it.

## Secrets

Create a real env file outside tracked paths:

```bash
cp deploy/kubernetes/base/secrets.example.env /tmp/twir-secrets.env
chmod 600 /tmp/twir-secrets.env
${EDITOR:-vi} /tmp/twir-secrets.env
```

Required keys:

```text
TWIR_DOPPLER_TOKEN
TWIR_POSTGRES_USER
TWIR_POSTGRES_DB
TWIR_POSTGRES_PASSWORD
CLICKHOUSE_USER
CLICKHOUSE_PASSWORD
CLICKHOUSE_DB
TEMPORAL_POSTGRES_USER
TEMPORAL_POSTGRES_PASSWORD
TEMPORAL_POSTGRES_DB
POSTGRES_BACKUP_S3_ENDPOINT
POSTGRES_BACKUP_S3_REGION
POSTGRES_BACKUP_S3_BUCKET
POSTGRES_BACKUP_S3_ACCESS_KEY_ID
POSTGRES_BACKUP_S3_SECRET_ACCESS_KEY
```

Create the namespace and secrets:

```bash
kubectl create namespace twir --dry-run=client -o yaml | kubectl apply -f -

kubectl -n twir create secret docker-registry registry-twir-app \
  --docker-server=registry.twir.app \
  --docker-username="${REGISTRY_USER}" \
  --docker-password="${REGISTRY_PASSWORD}" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl -n twir create secret generic twir-secrets \
  --from-env-file=/tmp/twir-secrets.env \
  --dry-run=client -o yaml | kubectl apply -f -
```

Do not commit `/tmp/twir-secrets.env` or any filled secret file.

## Render And Deploy

Preview the generated manifests:

```bash
kubectl kustomize deploy/kubernetes/overlays/prod
# or
kustomize build deploy/kubernetes/overlays/prod
```

Apply the full stack:

```bash
kubectl apply -k deploy/kubernetes/overlays/prod
```

Wait for stateful dependencies:

```bash
kubectl -n twir rollout status statefulset/postgres --timeout=15m
kubectl -n twir rollout status statefulset/clickhouse --timeout=15m
kubectl -n twir rollout status statefulset/redis --timeout=10m
kubectl -n twir rollout status statefulset/temporal-postgres --timeout=15m
kubectl -n twir rollout status statefulset/nats --timeout=10m
```

Wait for poolers:

```bash
kubectl -n twir rollout status deploy/pgdog --timeout=10m
kubectl -n twir rollout status deploy/pgbouncer --timeout=10m
```

Run or re-run migrations. The overlay creates `job/migrations`; delete and re-apply it when you need a fresh run:

```bash
kubectl -n twir delete job migrations --ignore-not-found
kubectl apply -f deploy/kubernetes/base/jobs.yaml
kubectl -n twir wait --for=condition=complete job/migrations --timeout=15m
kubectl -n twir logs job/migrations
```

If the migration job fails:

```bash
kubectl -n twir describe job migrations
kubectl -n twir logs job/migrations --all-containers=true
```

Wait for app workloads:

```bash
kubectl -n twir rollout status deploy/api-gql --timeout=10m
kubectl -n twir rollout status deploy/bots --timeout=10m
kubectl -n twir rollout status deploy/parser --timeout=10m
kubectl -n twir rollout status statefulset/eventsub --timeout=10m
kubectl -n twir rollout status deploy/web --timeout=10m
kubectl -n twir rollout status deploy/dashboard --timeout=10m
kubectl -n twir rollout status deploy/overlays --timeout=10m
kubectl -n twir rollout status deploy/websockets --timeout=10m
kubectl -n twir get pods -o wide
```

If apps started before migrations completed and entered CrashLoopBackOff, restart them after migrations pass:

```bash
kubectl -n twir rollout restart deploy/api-gql deploy/bots deploy/parser deploy/timers deploy/scheduler
kubectl -n twir rollout restart deploy/integrations deploy/web deploy/dashboard deploy/overlays deploy/websockets
kubectl -n twir rollout restart deploy/tokens deploy/emotes-cacher deploy/events deploy/language-processor deploy/toxicity-detector deploy/music-recognizer
kubectl -n twir rollout restart statefulset/eventsub
```

## Public Routes

The k3s Traefik replacement includes these hosts:

- `twir.app`
- `cf.twir.app`
- `services-bots.twir.app`
- `music-recognizer.twir.app`

Routes:

- `/api` -> `api-gql:3009`, with Traefik strip-prefix middleware in `IngressRoute`.
- `/s/` -> `api-gql:3009`, rewritten to `/v1/short-links/` by Traefik middleware.
- `/dashboard` -> `dashboard:8080`, with strip-prefix middleware.
- `/overlays` -> `overlays:8080`, with strip-prefix middleware.
- `/socket` -> `websockets:3004`, with strip-prefix middleware and Traefik `ServersTransport` timeout placeholders.
- `/` -> `web:3000`.
- `services-bots.twir.app` -> `bots:3000`.
- `music-recognizer.twir.app` -> `music-recognizer:3000`.

Check ingress after DNS or with `--resolve`:

```bash
kubectl -n twir get ingressroute,middleware,serverstransport
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/api/health
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/dashboard/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/overlays/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/s/test
curl -I --resolve services-bots.twir.app:443:NODE1_IP https://services-bots.twir.app/
curl -I --resolve music-recognizer.twir.app:443:NODE1_IP https://music-recognizer.twir.app/
```

For the exact HAProxy custom-domain fallback (`unknown host/path` -> `api-gql /v1/short-links{path}`), keep HAProxy or implement an explicit Traefik fallback route after deciding which custom domains are allowed.

## Manual Backup

The default Postgres backup CronJob runs daily at `03:00` cluster time. Override `spec.schedule` in an overlay if needed.

Run a manual backup from the CronJob:

```bash
kubectl -n twir create job --from=cronjob/postgres-backup postgres-backup-manual-$(date +%Y%m%d%H%M%S)
kubectl -n twir get jobs -l app.kubernetes.io/name=postgres-backup
kubectl -n twir logs -l app.kubernetes.io/name=postgres-backup --tail=200
```

A backup is not production-ready until you have restored it into a separate temporary Postgres and verified application data.

## Operational Checks

```bash
kubectl -n twir get all
kubectl -n twir get pvc
kubectl -n twir get pods -o wide
kubectl -n twir describe pvc postgres-data
kubectl -n twir logs deploy/otel-collector --tail=100
kubectl -n twir logs statefulset/postgres --tail=100
kubectl -n twir logs deploy/api-gql --tail=100
```

Postgres PVC migration note: compose mounted `/var/lib/postgresql`, while this manifest uses the official image's safer `/var/lib/postgresql/data` plus `PGDATA=/var/lib/postgresql/data/pgdata`. Do not mount an existing compose volume directly without planning and testing a data copy.
