# Развертывание Twir в Kubernetes на трех серверах Hetzner

## Цель

Этот документ описывает практический перенос Twir из `docker-compose.stack.yml` в Kubernetes на трех серверах Hetzner. Он покрывает кластер, хранение данных, секреты, ingress, порядок деплоя, проверку, резервные копии, откат и основные риски.

Рекомендованный вариант по умолчанию: `k3s` в HA режиме с embedded etcd на всех трех серверах.

Почему так:

- `k3s` проще полного Kubernetes, но остается обычным Kubernetes API.
- Embedded etcd на трех нодах дает кворум control plane без внешней базы.
- Три сервера достаточно для минимального HA control plane, если переживать отказ одной ноды.
- Не нужен managed Kubernetes. Все можно поднять на bare metal или обычных Hetzner VM.
- Rancher не обязателен. Rancher это UI, централизованное управление кластерами, политики и удобный day 2 ops. Для одного кластера Twir он полезен, но не нужен для работы k3s, деплоя приложений, ingress, storage, backup и мониторинга.

Важное ограничение: этот план не делает Postgres, ClickHouse, Redis, NATS JetStream и Temporal Postgres полностью HA. С local PV они привязаны к конкретной ноде. Для HA баз нужны отдельная репликация, оператор или внешний managed сервис.


## Готовые Kubernetes manifests

В репозитории добавлен deployable kustomize layout: `deploy/kubernetes`.

Быстрый production entrypoint:

```bash
kubectl kustomize deploy/kubernetes/overlays/prod
kubectl apply -k deploy/kubernetes/overlays/prod
```

Подробные команды для k3s, labels, registry pull secret, `twir-secrets`, deploy order, migrations, ingress checks и manual backup находятся в `deploy/kubernetes/README.md`.

Минимальный quickstart после установки k3s:

```bash
kubectl label node twir-k3s-1 twir.app/database-primary=true twir.app/database=true --overwrite
kubectl label node twir-k3s-2 twir.app/database=true twir.app/cache=true twir.app/backup=true --overwrite
kubectl label node twir-k3s-3 twir.app/apps=true --overwrite

kubectl create namespace twir --dry-run=client -o yaml | kubectl apply -f -

kubectl -n twir create secret docker-registry registry-twir-app \
  --docker-server=registry.twir.app \
  --docker-username="${REGISTRY_USER}" \
  --docker-password="${REGISTRY_PASSWORD}" \
  --dry-run=client -o yaml | kubectl apply -f -

cp deploy/kubernetes/base/secrets.example.env /tmp/twir-secrets.env
chmod 600 /tmp/twir-secrets.env
${EDITOR:-vi} /tmp/twir-secrets.env

kubectl -n twir create secret generic twir-secrets \
  --from-env-file=/tmp/twir-secrets.env \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -k deploy/kubernetes/overlays/prod

kubectl -n twir rollout status statefulset/postgres --timeout=15m
kubectl -n twir rollout status statefulset/clickhouse --timeout=15m
kubectl -n twir rollout status statefulset/redis --timeout=10m
kubectl -n twir rollout status statefulset/temporal-postgres --timeout=15m
kubectl -n twir rollout status statefulset/nats --timeout=10m

kubectl -n twir wait --for=condition=complete job/migrations --timeout=15m
kubectl -n twir logs job/migrations

kubectl -n twir rollout status deploy/api-gql --timeout=10m
kubectl -n twir rollout status deploy/web --timeout=10m
kubectl -n twir rollout status deploy/dashboard --timeout=10m
kubectl -n twir rollout status deploy/websockets --timeout=10m
kubectl -n twir get pods -o wide
```

Важные ограничения manifests:

- Они не делают DB HA. Stateful базы используют single replica и `local-path` PVC.
- HAProxy из compose заменен k3s Traefik Ingress/IngressRoute. Для точного custom-domain fallback, HAProxy cache, gzip и stick-table rate limit оставьте HAProxy отдельным edge DaemonSet или внешним proxy.
- Перед production cutover замените `latest` tags на immutable tags или image digests.
- Real secrets создаются только через Kubernetes Secret из env file; не заполняйте `secret-template.yaml` в git.

## Предположения и значения для замены

Замените эти значения перед запуском:

```text
NODE1_IP=<public-ip-1>
NODE2_IP=<public-ip-2>
NODE3_IP=<public-ip-3>
NODE1_HOST=twir-k3s-1
NODE2_HOST=twir-k3s-2
NODE3_HOST=twir-k3s-3
K3S_TOKEN=<long-random-token>
REGISTRY_USER=<registry.twir.app-user>
REGISTRY_PASSWORD=<registry.twir.app-password>
TWIR_DOPPLER_TOKEN=<doppler-token>
TWIR_POSTGRES_USER=<postgres-user>
TWIR_POSTGRES_DB=<postgres-db>
TWIR_POSTGRES_PASSWORD=<postgres-password>
```

Пример роли нод:

- `twir-k3s-1`: DB primary node, label `twir.app/database-primary=true`, аналог Swarm label `databases-1 == true`.
- `twir-k3s-2`: DB, backup, cache node, label `twir.app/database=true`, аналог Swarm label `databases == true`.
- `twir-k3s-3`: app worker node.

Так как серверов всего три, control plane и workload будут делить CPU, RAM, disk IO и сеть. Задайте `requests` и `limits`, иначе приложения могут вытеснить базы или etcd.

## Архитектура кластера

Минимальная схема:

```text
Internet
  |
DNS A/AAAA: twir.app, cf.twir.app, services-bots.twir.app, music-recognizer.twir.app
  |
3 Hetzner servers, ports 80/443 open
  |
k3s HA, embedded etcd on all nodes
  |
Ingress Controller, Traefik from k3s or nginx-ingress
  |
Namespace twir
  |
StatefulSets, Deployments, Services, Jobs, ConfigMaps, Secrets, PVCs
```

Компоненты из Swarm:

- Edge routing: HAProxy `haproxy:3.3-alpine`, global, host ports `80` и `443`, config `configs/haproxy/haproxy.cfg`.
- Stateful dependencies: Postgres 18, ClickHouse, Redis 8, Temporal Postgres, NATS JetStream.
- Poolers: PgDog 3 replicas, PgBouncer 5 replicas.
- App workloads: API, bots, parser, timers, scheduler, eventsub, integrations, web, dashboard, overlays, websockets, tokens, emotes-cacher, events, tts, language-processor, toxicity-detector, music-recognizer.
- Ops workloads: migrations, postgres-backup, otel-collector, adminer, temporal, temporal-ui.

## Подготовка серверов

Выполните на всех трех серверах Ubuntu или Debian-like.

```bash
sudo hostnamectl set-hostname twir-k3s-1
```

На второй и третьей ноде используйте `twir-k3s-2` и `twir-k3s-3`.

```bash
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install -y curl ca-certificates gnupg lsb-release jq open-iscsi nfs-common
```

Отключите swap:

```bash
sudo swapoff -a
sudo sed -i.bak '/ swap / s/^/#/' /etc/fstab
```

Модули ядра:

```bash
cat <<'EOF' | sudo tee /etc/modules-load.d/k8s.conf
br_netfilter
overlay
EOF
sudo modprobe br_netfilter
sudo modprobe overlay
```

Sysctl:

```bash
cat <<'EOF' | sudo tee /etc/sysctl.d/99-kubernetes.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward = 1
vm.max_map_count = 262144
fs.file-max = 2097152
EOF
sudo sysctl --system
```

Firewall. Если используете `ufw`, откройте минимум:

```bash
sudo ufw allow 22/tcp
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw allow 6443/tcp
sudo ufw allow 2379:2380/tcp
sudo ufw allow 8472/udp
sudo ufw allow 10250/tcp
sudo ufw allow from NODE1_IP
sudo ufw allow from NODE2_IP
sudo ufw allow from NODE3_IP
sudo ufw enable
```

Для Hetzner Cloud Firewall лучше ограничить `6443`, `2379-2380`, `8472/udp`, `10250` только IP трех серверов. Порты `80/443` оставьте публичными.

## Установка k3s HA с embedded etcd

### Первая нода

На `twir-k3s-1`:

```bash
curl -sfL https://get.k3s.io | \
  K3S_TOKEN="${K3S_TOKEN}" \
  sh -s - server \
  --cluster-init \
  --tls-san "${NODE1_IP}" \
  --tls-san "${NODE2_IP}" \
  --tls-san "${NODE3_IP}" \
  --write-kubeconfig-mode 644
```

Проверьте:

```bash
sudo kubectl get nodes -o wide
sudo kubectl get pods -A
```

### Вторая и третья ноды

На `twir-k3s-2`:

```bash
curl -sfL https://get.k3s.io | \
  K3S_TOKEN="${K3S_TOKEN}" \
  sh -s - server \
  --server https://${NODE1_IP}:6443 \
  --tls-san "${NODE1_IP}" \
  --tls-san "${NODE2_IP}" \
  --tls-san "${NODE3_IP}" \
  --write-kubeconfig-mode 644
```

На `twir-k3s-3` повторите ту же команду.

### Kubeconfig для локальной машины

На первой ноде:

```bash
sudo cat /etc/rancher/k3s/k3s.yaml
```

Скопируйте файл локально в `~/.kube/twir-k3s.yaml` и замените `https://127.0.0.1:6443` на `https://NODE1_IP:6443` или на DNS/VIP API endpoint.

```bash
export KUBECONFIG=~/.kube/twir-k3s.yaml
kubectl get nodes -o wide
```

Для production лучше держать API server за private IP или VPN, не за публичным адресом.

## Labels и taints

Назначьте labels:

```bash
kubectl label node twir-k3s-1 twir.app/database-primary=true twir.app/database=true
kubectl label node twir-k3s-2 twir.app/database=true twir.app/backup=true twir.app/cache=true
kubectl label node twir-k3s-3 twir.app/apps=true
```

Для app workload можно использовать `nodeAffinity`, чтобы не класть приложения на DB-ноды. Но в кластере из трех серверов жесткий запрет может привести к нехватке capacity. Практичный вариант:

- Базы закрепить `requiredDuringSchedulingIgnoredDuringExecution` на нужных DB labels.
- Приложениям дать `preferredDuringSchedulingIgnoredDuringExecution` против `twir.app/database=true`, а не жесткий запрет.
- Worker-only сервисам `language-processor` и `toxicity-detector` дать `required` на `twir.app/apps=true`, если третья нода достаточно большая.

Пример affinity для Postgres:

```yaml
nodeSelector:
  twir.app/database-primary: "true"
```

Пример мягкого избегания DB нод для приложений:

```yaml
affinity:
  nodeAffinity:
    preferredDuringSchedulingIgnoredDuringExecution:
      - weight: 100
        preference:
          matchExpressions:
            - key: twir.app/database
              operator: DoesNotExist
```

## Как Swarm переводится в Kubernetes

| Swarm | Kubernetes |
| --- | --- |
| `deploy.replicas` | `Deployment.spec.replicas` или `StatefulSet.spec.replicas` |
| `deploy.mode: global` | `DaemonSet`, для HAProxy на каждом node |
| `placement.constraints` | `nodeSelector`, `nodeAffinity`, `podAntiAffinity`, taints/tolerations |
| `secrets` | `Secret`, mount file или env var |
| `configs` | `ConfigMap`, mount file |
| named volumes | `PersistentVolumeClaim` и `PersistentVolume` |
| external networks | Kubernetes `Service`, `Ingress`, NetworkPolicy при необходимости |
| `endpoint_mode: dnsrr` | обычный Kubernetes Service DNS, headless Service если нужен прямой DNS pod endpoints |
| `update_config.order: start-first` | rolling update с `maxUnavailable: 0`, `maxSurge: 1` |
| `update_config.order: stop-first` | rolling update с `maxSurge: 0`, `maxUnavailable: 1`, или ручной rollout |
| host published ports | Ingress Controller Service, `hostNetwork`, `hostPort` или DaemonSet с host ports |
| `depends_on` | readiness/startup probes, initContainers, Jobs, deploy order в GitOps или CI |
| restart policy | `restartPolicy: Always` для Deployments/StatefulSets, `backoffLimit` для Jobs |

Kubernetes не гарантирует порядок старта через `depends_on`. Готовность нужно выражать readiness probes, initContainers, Jobs и порядком применения манифестов.

## Namespace, registry и базовые Secrets

Создайте namespace:

```bash
kubectl create namespace twir
```

Registry pull secret для `registry.twir.app`:

```bash
kubectl -n twir create secret docker-registry registry-twir-app \
  --docker-server=registry.twir.app \
  --docker-username="${REGISTRY_USER}" \
  --docker-password="${REGISTRY_PASSWORD}"
```

Секреты приложения:

```bash
kubectl -n twir create secret generic twir-secrets \
  --from-literal=twir_doppler_token="${TWIR_DOPPLER_TOKEN}" \
  --from-literal=twir_postgres_user="${TWIR_POSTGRES_USER}" \
  --from-literal=twir_postgres_db="${TWIR_POSTGRES_DB}" \
  --from-literal=twir_postgres_password="${TWIR_POSTGRES_PASSWORD}"
```

Если контейнеры ожидают файлы как в Swarm, монтируйте Secret в `/run/secrets` с именами ключей:

```yaml
volumes:
  - name: twir-secrets
    secret:
      secretName: twir-secrets
containers:
  - name: postgres
    volumeMounts:
      - name: twir-secrets
        mountPath: /run/secrets
        readOnly: true
```

Для Doppler token чаще проще использовать env:

```yaml
env:
  - name: DOPPLER_TOKEN
    valueFrom:
      secretKeyRef:
        name: twir-secrets
        key: twir_doppler_token
```

Проверьте точное имя env var, которое ожидает каждый образ. В compose секрет называется `twir_doppler_token`, но образ может читать файл или env.

## ConfigMaps

### ClickHouse

Создайте ConfigMap из текущего XML:

```bash
kubectl -n twir create configmap clickhouse-config \
  --from-file=override.xml=configs/clickhouse-config.xml
```

Смысл текущего конфига:

- Отключены тяжелые system logs: metric, query_thread, part, trace, crash, OpenTelemetry span и другие.
- `query_log` оставлен с TTL 1 день.
- Logger выставлен в `warning`.
- `log_queries` и `log_query_threads` выключены в default profile.

Монтирование:

```yaml
volumeMounts:
  - name: clickhouse-config
    mountPath: /etc/clickhouse-server/config.d/override.xml
    subPath: override.xml
volumes:
  - name: clickhouse-config
    configMap:
      name: clickhouse-config
```

### PgDog

Текущий `configs/pgdog/pgdog.prod.toml`:

- `min_pool_size = 1`
- `default_pool_size = 30`
- `pooler_mode = "transaction"`
- database host `postgres`, port `5432`

`configs/pgdog/users.toml` сейчас содержит статический пример `twir/twir`. В Kubernetes это нельзя хранить как обычный ConfigMap. Сделайте шаблон через Secret или initContainer, который собирает `users.toml` из `twir_postgres_user` и `twir_postgres_password`.

Минимальный вариант:

```bash
kubectl -n twir create configmap pgdog-config \
  --from-file=pgdog.toml=configs/pgdog/pgdog.prod.toml \
  --from-file=start.sh=configs/pgdog/start.sh
```

А `users.toml` создавайте Secret:

```bash
kubectl -n twir create secret generic pgdog-users \
  --from-literal=users.toml="$(cat <<EOF
[[users]]
name = \"${TWIR_POSTGRES_USER}\"
database = \"${TWIR_POSTGRES_DB}\"
password = \"${TWIR_POSTGRES_PASSWORD}\"
server_user = \"${TWIR_POSTGRES_USER}\"
server_password = \"${TWIR_POSTGRES_PASSWORD}\"
EOF
)"
```

## Стратегия хранения на Hetzner

### Рекомендация по умолчанию

Для Postgres, ClickHouse, Redis, Temporal Postgres и NATS JetStream используйте local PV или стандартный `local-path` k3s для лучшей latency и predictable IO.

Это не HA для данных. Если нода с local PV умерла, pod не переедет на другую ноду с теми же данными. Нужно восстановление из backup или отдельная репликация.

Практический вариант:

- Postgres 18, ClickHouse, Temporal Postgres закрепить на `twir-k3s-1` через `twir.app/database-primary=true`.
- Redis и postgres-backup закрепить на `twir-k3s-2` через `twir.app/database=true` или `twir.app/cache=true`.
- NATS JetStream можно начать с одного StatefulSet replica на app node или DB node, но для JetStream HA нужен кластер NATS с PVC на разных нодах и quorum JetStream.

### Longhorn

Longhorn можно поставить, если нужен replicated block storage между серверами. Но для баз данных это компромисс:

- Плюс: pod может переехать на другую ноду с томом.
- Минус: network replication увеличивает latency и write amplification.
- Риск: при слабой сети и перегрузе дисков базы будут страдать.

Для Postgres и ClickHouse на трех Hetzner серверах чаще лучше local PV плюс хорошие backup, чем Longhorn без тщательного теста нагрузки.

### Hetzner Cloud CSI

Hetzner Cloud CSI используйте только если эти серверы являются Hetzner Cloud VMs в одной location, где доступны attachable volumes. CSI volume обычно может быть attached к одной VM за раз. Это удобно для переезда pod, но не делает базу HA и может иметь latency хуже локального NVMe.

### Пример PVC

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-data
  namespace: twir
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: local-path
  resources:
    requests:
      storage: 200Gi
```

Для production лучше явно создать local PV на конкретном диске, например `/var/lib/twir/postgres`, чтобы не потерять данные при очистке системных директорий.

## Stateful services

### Postgres 18

Swarm использует `postgres:18.1` и tuned command:

```text
postgres -c max_connections=1000 -c shared_buffers=2GB -c effective_cache_size=6GB -c maintenance_work_mem=512MB -c checkpoint_completion_target=0.9 -c wal_buffers=16MB -c default_statistics_target=100 -c random_page_cost=1.1 -c effective_io_concurrency=200 -c work_mem=32MB -c huge_pages=off -c min_wal_size=1GB -c max_wal_size=4GB -c shared_preload_libraries='pg_stat_statements'
```

Kubernetes ресурс: `StatefulSet` с 1 replica, `Service` `postgres`, PVC `postgres-data`, `nodeSelector twir.app/database-primary=true`.

Сохраняйте secret file env из compose:

```yaml
env:
  - name: POSTGRES_USER_FILE
    value: /run/secrets/twir_postgres_user
  - name: POSTGRES_PASSWORD_FILE
    value: /run/secrets/twir_postgres_password
  - name: POSTGRES_DB_FILE
    value: /run/secrets/twir_postgres_db
```

Volume path должен быть `/var/lib/postgresql`, как в compose. Проверьте image docs Postgres 18. Если image ожидает `/var/lib/postgresql/data`, скорректируйте PVC mount и `PGDATA`. Не переносите данные вслепую.

### ClickHouse

Swarm использует `clickhouse/clickhouse-server:25.5-alpine`, 1 replica, limits `2 CPU`, `3G`, PVC `/var/lib/clickhouse`, ConfigMap XML в `/etc/clickhouse-server/config.d/override.xml`, env `CLICKHOUSE_USER`, `CLICKHOUSE_PASSWORD`, `CLICKHOUSE_DB`, `CLICKHOUSE_DEFAULT_ACCESS_MANAGEMENT=1`.

Kubernetes ресурс: `StatefulSet` с 1 replica, `Service` `clickhouse`, PVC `clickhouse-data`, ConfigMap `clickhouse-config`, nodeSelector `twir.app/database-primary=true`.

Перенесите `CLICKHOUSE_PASSWORD` в Secret. Значение `twir` из compose годится только как пример.

### Redis 8

Swarm использует `redis:8.0.2` и команду:

```text
redis-server --save 60 1 --loglevel warning --io-threads 4
```

Kubernetes ресурс: `StatefulSet` с 1 replica, `Service` `redis`, PVC `redis-data`, nodeSelector `twir.app/cache=true` или `twir.app/database=true`.

### Temporal Postgres

Swarm использует `bitnamilegacy/postgresql:17`, database `temporal`, volume `/bitnami/postgresql`, label `databases-1 == true`.

Kubernetes ресурс: `StatefulSet` с 1 replica, `Service` `temporal-postgres`, PVC `temporal-postgres-data`, nodeSelector `twir.app/database-primary=true`.

Вынесите `POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRESQL_POSTGRES_PASSWORD` в Secret. Значения `temporal/temporal` из compose замените.

### NATS JetStream

Swarm использует `nats:2.10.11-scratch`, command `-js -m 8222`, без named volume. Для production JetStream без PVC рискован.

Kubernetes ресурс минимум: `StatefulSet` с 1 replica, `Service` `nats` на `4222`, metrics `8222`, PVC для JetStream store.

Пример command:

```yaml
args:
  - -js
  - -m
  - "8222"
  - -sd
  - /data/jetstream
```

Для HA JetStream позже делайте 3 replica NATS cluster с pod anti-affinity и PVC на разных нодах. Это отдельная миграция, не просто увеличение replicas.

## Poolers

### PgDog

Swarm: `ghcr.io/pgdogdev/pgdog:v0.1.34`, 3 replicas, transaction pooling, `default_pool_size=30`, `min_pool_size=1`, database host `postgres:5432`, update `stop-first` parallelism 1.

Kubernetes ресурс: `Deployment` 3 replicas, `Service` `pgdog`, rolling update ближе к stop-first:

```yaml
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxSurge: 0
    maxUnavailable: 1
```

### PgBouncer

Swarm: `bitnamilegacy/pgbouncer:1.23.1`, 5 replicas, env:

```text
POSTGRESQL_HOST=postgres
PGBOUNCER_AUTH_USER=twir
PGBOUNCER_DATABASE=twir
PGBOUNCER_DEFAULT_POOL_SIZE=19
PGBOUNCER_MIN_POOL_SIZE=10
PGBOUNCER_MAX_CLIENT_CONN=4000
PGBOUNCER_POOL_MODE=transaction
```

Kubernetes ресурс: `Deployment` 5 replicas, `Service` `pgbouncer`, rolling update `maxUnavailable: 0`, `maxSurge: 1` или `2`.

Проверьте, что суммарные pool sizes PgDog и PgBouncer не превышают реальный `max_connections=1000` Postgres с запасом для админских соединений, migrations и backup.

## Ingress, TLS и маршрутизация

### Рекомендация

Для Kubernetes начните с встроенного Traefik из k3s или поставьте nginx-ingress. Это ближе к обычному Kubernetes и проще поддерживается через `Ingress` или Traefik CRD `IngressRoute`.

HAProxy из compose можно оставить как `DaemonSet` с `hostNetwork: true` или `hostPort: 80/443`, если вам нужно сохранить точное поведение HAProxy:

- stick-table rate limit `1000 req / 10s` по `X-Ru-Detected-IP` или source IP;
- HAProxy cache `mycache` для dashboard и web;
- gzip compression прямо на edge;
- custom-domain fallback через default backend;
- долгий WebSocket tunnel timeout `3600s`.

Если точное совпадение не критично, используйте ingress controller и перенесите часть логики в middleware, annotations или app level.

### Что делает HAProxy сейчас

Из `configs/haproxy/haproxy.cfg`:

- Image `haproxy:3.3-alpine`.
- Swarm `deploy.mode: global`.
- Host ports `80` и `443`.
- Config target `/usr/local/etc/haproxy/haproxy.cfg`.
- Hosts: `twir.app`, `cf.twir.app`, `services-bots.twir.app`, `music-recognizer.twir.app`.
- `/api` идет в `api-gql:3009`, prefix `/api` снимается.
- `/s/` переписывается в `/v1/short-links/` и идет в `api-gql`.
- `/dashboard` идет в `dashboard:8080`, prefix снимается.
- `/overlays` идет в `overlays:8080`, prefix снимается.
- `/socket` идет в `websockets:3004`, prefix снимается, tunnel timeout `3600s`.
- `services-bots.twir.app` идет в `bots:3000`.
- `music-recognizer.twir.app` идет в `music-recognizer:3000`.
- Остальные hosts и paths идут в fallback `api-gql` как `/v1/short-links{path}`.
- Rate limit: 1000 requests за 10 секунд, ключ `X-Ru-Detected-IP`, если header есть, иначе source IP.
- Gzip compression для JSON, GraphQL response, JS, CSS, HTML, plain text.
- Маленький cache для dashboard и web: max age 240 sec, max object 10000 bytes, total 4 MB.

### Service ports

Создайте ClusterIP services:

```text
api-gql: 3009
bots: 3000
music-recognizer: 3000
web: 3000
dashboard: 8080
overlays: 8080
websockets: 3004
temporal: 7233
temporal-ui: 8080
adminer: 8080
nats: 4222, 8222
postgres: 5432
pgdog: проверить порт образа PgDog
pgbouncer: 6432
redis: 6379
clickhouse: 8123, 9000
```

Порт PgDog проверьте по образу и текущему `start.sh`, так как compose не публикует порт наружу.

### Traefik IngressRoute пример

Если используете Traefik CRD, для strip prefix:

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: api-strip
  namespace: twir
spec:
  stripPrefix:
    prefixes:
      - /api
---
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: socket-strip
  namespace: twir
spec:
  stripPrefix:
    prefixes:
      - /socket
---
apiVersion: traefik.io/v1alpha1
kind: IngressRoute
metadata:
  name: twir-main
  namespace: twir
spec:
  entryPoints:
    - websecure
  routes:
    - match: Host(`twir.app`) && PathPrefix(`/api`)
      kind: Rule
      middlewares:
        - name: api-strip
      services:
        - name: api-gql
          port: 3009
    - match: Host(`twir.app`) && PathPrefix(`/socket`)
      kind: Rule
      middlewares:
        - name: socket-strip
      services:
        - name: websockets
          port: 3004
    - match: Host(`twir.app`)
      kind: Rule
      services:
        - name: web
          port: 3000
  tls:
    secretName: twir-app-tls
```

Для `/s/` rewrite в Traefik нужен `replacePathRegex`:

```yaml
apiVersion: traefik.io/v1alpha1
kind: Middleware
metadata:
  name: short-links-rewrite
  namespace: twir
spec:
  replacePathRegex:
    regex: ^/s/(.*)
    replacement: /v1/short-links/$1
```

Для nginx-ingress используйте отдельные Ingress objects и annotations `nginx.ingress.kubernetes.io/rewrite-target`, `nginx.ingress.kubernetes.io/proxy-read-timeout: "3600"`, `nginx.ingress.kubernetes.io/proxy-send-timeout: "3600"`.

### TLS

Поставьте cert-manager или используйте статические TLS Secrets.

cert-manager пример:

```bash
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.16.2/cert-manager.yaml
```

После установки создайте `ClusterIssuer` Let's Encrypt. Если DNS находится за Cloudflare proxy, проверьте real IP headers. Текущий HAProxy rate limit использует `X-Ru-Detected-IP`, такой header нужно сохранить или заменить на ingress-real-ip настройки.

### DNS cutover

До миграции:

- Уменьшите TTL DNS до 60 или 300 секунд за сутки.
- Поднимите ingress на временном hostname, например `k8s-preview.twir.app`.
- Протестируйте routes через `curl --resolve`.
- Переключите `A/AAAA` записей на IP нод ingress или на внешний load balancer.

Если нет внешнего load balancer, можно указать несколько A records на все три сервера. Это проще, но health checking будет слабее, чем у LB.

## Порядок деплоя

Рекомендуемый порядок:

1. Namespace, labels, StorageClass, PV/PVC, image pull secret.
2. Secrets и ConfigMaps.
3. Stateful dependencies: Postgres, ClickHouse, Redis, Temporal Postgres, NATS JetStream.
4. Poolers: PgDog, PgBouncer.
5. Temporal service и Temporal UI.
6. Migrations Job.
7. Internal app services: parser, timers, scheduler, eventsub, integrations, tokens, emotes-cacher, events, tts, language-processor, toxicity-detector.
8. Public app services: api-gql, bots, music-recognizer, web, dashboard, overlays, websockets.
9. Ingress, TLS, DNS preview.
10. Observability: otel-collector, logs, metrics, alerts.
11. Backups: postgres-backup, ClickHouse backups, Redis/NATS policy.
12. DNS cutover.

Не запускайте app workloads до успешных migrations. Не переключайте production DNS до проверки backup restore rehearsal.

## Таблица компонентов

| Compose service | Image | Kubernetes resource | Replicas | Placement | Notes |
| --- | --- | --- | ---: | --- | --- |
| `haproxy` | `haproxy:3.3-alpine` | Лучше Ingress Controller. Если сохранять HAProxy, `DaemonSet` | по нодам | all nodes, host ports 80/443 | ConfigMap from `configs/haproxy/haproxy.cfg`; exact rate limit/cache/rewrite behavior |
| `adminer` | `adminer` | `Deployment` + `Service` | 1 | avoid DB nodes | Лучше закрыть за auth/VPN, не публиковать публично |
| `nats` | `nats:2.10.11-scratch` | `StatefulSet` + PVC + `Service` | 1 initially | app or DB node | Add PVC for JetStream, compose has no named volume |
| `otel-collector` | `otel/opentelemetry-collector-contrib:0.116.1` | `DaemonSet` или `Deployment` | 1 или per node | control plane preferred | Mount config from ConfigMap, Docker socket заменить на Kubernetes receiver или node logs pipeline |
| `postgres` | `postgres:18.1` | `StatefulSet` + `Service` | 1 | `twir.app/database-primary=true` | Keep tuned command and file secrets |
| `postgres-backup` | `registry.twir.app/twirapp/postgres-backup:latest` | `CronJob` или `Deployment` | 1 | `twir.app/backup=true` | Better as CronJob if image supports one-shot backup |
| `pgdog` | `ghcr.io/pgdogdev/pgdog:v0.1.34` | `Deployment` + `Service` | 3 | apps preferred | `stop-first` style rolling update, users as Secret |
| `pgbouncer` | `bitnamilegacy/pgbouncer:1.23.1` | `Deployment` + `Service` | 5 | apps preferred | Transaction pooling, env from Secret |
| `clickhouse` | `clickhouse/clickhouse-server:25.5-alpine` | `StatefulSet` + `Service` | 1 | `twir.app/database-primary=true` | ConfigMap XML, Secret password, PVC |
| `temporal-postgres` | `bitnamilegacy/postgresql:17` | `StatefulSet` + `Service` | 1 | `twir.app/database-primary=true` | Replace static credentials |
| `temporal` | `twirapp/temporal:latest` | `Deployment` + `Service` | 1 | avoid DB nodes | Depends on DB readiness, pin image tag |
| `temporal-ui` | `temporalio/ui:2.21.0` | `Deployment` + `Service` | 1 | avoid DB nodes | `TEMPORAL_ADDRESS=temporal:7233` |
| `migrations` | `registry.twir.app/twirapp/migrations:latest` | `Job` | 1 run | apps preferred | Run after DB and before apps |
| `redis` | `redis:8.0.2` | `StatefulSet` + `Service` | 1 | `twir.app/cache=true` | Command from compose, PVC |
| `api-gql` | `registry.twir.app/twirapp/api-gql:latest` | `Deployment` + `Service` | 1 | avoid DB nodes | Port 3009, public via ingress |
| `bots` | `registry.twir.app/twirapp/bots:latest` | `Deployment` + `Service` | 4 | avoid DB nodes | Host `services-bots.twir.app`, port 3000 |
| `parser` | `registry.twir.app/twirapp/parser:latest` | `Deployment` | 6 | avoid DB nodes | Internal worker |
| `timers` | `registry.twir.app/twirapp/timers:latest` | `Deployment` | 1 | avoid DB nodes | Internal worker |
| `scheduler` | `registry.twir.app/twirapp/scheduler:latest` | `Deployment` | 1 | avoid DB nodes | Internal worker |
| `eventsub` | `registry.twir.app/twirapp/eventsub:latest` | `StatefulSet` preferred | 3 | avoid DB nodes | Compose uses `REPLICA={{.Task.Slot}}`; use StatefulSet ordinal or Downward API |
| `integrations` | `registry.twir.app/twirapp/integrations:latest` | `Deployment` | 1 | apps preferred | Internal service |
| `web` | `registry.twir.app/twirapp/web:latest` | `Deployment` + `Service` | 4 | avoid DB nodes | Port 3000, cache/gzip at ingress if needed |
| `dashboard` | `registry.twir.app/twirapp/dashboard:latest` | `Deployment` + `Service` | 1 | avoid DB nodes | Command from compose, port 8080, `/dashboard` strip |
| `overlays` | `registry.twir.app/twirapp/overlays:latest` | `Deployment` + `Service` | 1 | avoid DB nodes | Command from compose, port 8080, `/overlays` strip |
| `websockets` | `registry.twir.app/twirapp/websockets:latest` | `Deployment` + `Service` | 1 | avoid DB nodes | Port 3004, long timeout `/socket` |
| `tokens` | `registry.twir.app/twirapp/tokens:latest` | `Deployment` | 4 | avoid DB nodes | Internal worker/service |
| `emotes-cacher` | `registry.twir.app/twirapp/emotes-cacher:latest` | `Deployment` | 1 | avoid DB nodes | Internal worker |
| `events` | `registry.twir.app/twirapp/events:latest` | `Deployment` | 6 | avoid DB nodes | Internal worker |
| `tts` | `aculeasis/rhvoice-rest:latest` | `Deployment` + `Service` | 4 | avoid DB nodes | Internal HTTP service likely needed by apps |
| `language-processor` | `twirapp/language-processor-py:latest` | `Deployment` | 4 | worker-only | Compose excludes managers; use `twir.app/apps=true` if capacity allows |
| `toxicity-detector` | `registry.twir.app/twirapp/toxicity-detector:latest` | `Deployment` | 4 | worker-only | `TOXICITY_THRESHOLD=-4`, use app node if capacity allows |
| `music-recognizer` | `registry.twir.app/twirapp/music-recognizer:latest` | `Deployment` + `Service` | 4 | avoid DB nodes | Host `music-recognizer.twir.app`, port 3000 |

Предупреждение: большинство Twir images в compose используют `latest`. Для production закрепите immutable tags или digest, например `registry.twir.app/twirapp/api-gql:2026-06-01-abcdef` или `@sha256:...`. Без этого rollback может поднять уже другой образ с тем же tag.

## Скелеты манифестов

### Deployment app service

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gql
  namespace: twir
spec:
  replicas: 1
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxUnavailable: 0
      maxSurge: 1
  selector:
    matchLabels:
      app: api-gql
  template:
    metadata:
      labels:
        app: api-gql
    spec:
      imagePullSecrets:
        - name: registry-twir-app
      affinity:
        nodeAffinity:
          preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 100
              preference:
                matchExpressions:
                  - key: twir.app/database
                    operator: DoesNotExist
      containers:
        - name: api-gql
          image: registry.twir.app/twirapp/api-gql:REPLACE_WITH_PINNED_TAG
          ports:
            - containerPort: 3009
          env:
            - name: DOPPLER_TOKEN
              valueFrom:
                secretKeyRef:
                  name: twir-secrets
                  key: twir_doppler_token
          resources:
            requests:
              cpu: 500m
              memory: 512Mi
            limits:
              cpu: "2"
              memory: 3Gi
          readinessProbe:
            httpGet:
              path: /health
              port: 3009
            initialDelaySeconds: 10
            periodSeconds: 10
```

Если у сервиса нет `/health`, замените probe на рабочий endpoint или TCP probe. Не оставляйте нерабочий health check.

### Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: api-gql
  namespace: twir
spec:
  selector:
    app: api-gql
  ports:
    - name: http
      port: 3009
      targetPort: 3009
```

### Migrations Job

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: migrations
  namespace: twir
spec:
  backoffLimit: 3
  template:
    spec:
      restartPolicy: OnFailure
      imagePullSecrets:
        - name: registry-twir-app
      containers:
        - name: migrations
          image: registry.twir.app/twirapp/migrations:REPLACE_WITH_PINNED_TAG
          env:
            - name: DOPPLER_TOKEN
              valueFrom:
                secretKeyRef:
                  name: twir-secrets
                  key: twir_doppler_token
```

Перед повторным запуском удаляйте старый Job или используйте новое имя:

```bash
kubectl -n twir delete job migrations
kubectl -n twir apply -f migrations-job.yaml
```

## EventSub replica equivalent

Compose задает:

```yaml
environment:
  REPLICA: "{{.Task.Slot}}"
```

В Kubernetes у Deployment нет стабильного slot number. Используйте `StatefulSet` для `eventsub`, тогда pods будут `eventsub-0`, `eventsub-1`, `eventsub-2`.

Если приложение принимает строковый `REPLICA`, можно передать имя pod:

```yaml
env:
  - name: POD_NAME
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
  - name: REPLICA
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
```

Если нужен именно номер, добавьте entrypoint wrapper, который достает ordinal из hostname:

```sh
export REPLICA="${HOSTNAME##*-}"
exec /app/eventsub
```

Проверьте реальный command в image перед заменой entrypoint.

## Observability и логи

Минимум:

- `otel-collector` как `Deployment` или `DaemonSet` с ConfigMap из `configs/otel/otel-collector.yaml`.
- Логи pod через `kubectl logs`, затем Loki/Promtail или OpenTelemetry logs pipeline.
- Метрики ingress controller, k3s, node exporter, Postgres exporter, Redis exporter, ClickHouse exporter, NATS metrics `:8222`.
- Alerts на disk usage PVC, etcd health, Postgres connections, backup age, pod crash loops, 5xx ingress, latency WebSocket.

Docker socket mount из compose не переносите напрямую. В Kubernetes используйте Kubernetes metadata receiver, filelog receiver или DaemonSet collector с доступом к `/var/log/pods`.

Проверка OTel:

```bash
kubectl -n twir logs deploy/otel-collector
kubectl -n twir port-forward deploy/otel-collector 8888:8888
```

## Backup и restore

### Postgres

`postgres-backup` из compose лучше оформить как `CronJob`, если образ делает один backup и завершается. Если это daemon, оставьте `Deployment` 1 replica на node `twir.app/backup=true`.

Минимальный план:

- Ежедневный full backup через `pg_dump` или physical backup tool.
- WAL archive, если нужен point-in-time recovery.
- Хранение вне кластера: S3-compatible storage, Hetzner Storage Box, Backblaze, MinIO вне этих нод.
- Retention, например 7 daily, 4 weekly, 6 monthly.
- Шифрование backup и отдельный Secret для credentials.

Restore rehearsal обязателен:

```bash
kubectl -n twir scale deploy/api-gql --replicas=0
kubectl -n twir scale deploy/parser --replicas=0
kubectl -n twir run pg-restore-test --rm -it --image=postgres:18.1 -- bash
```

Внутри тестового pod восстановите backup в отдельную временную БД или отдельный временный Postgres StatefulSet. Не тестируйте restore поверх production данных.

Проверки после restore:

```bash
psql -h postgres -U "$TWIR_POSTGRES_USER" -d "$TWIR_POSTGRES_DB" -c 'select count(*) from information_schema.tables;'
```

### ClickHouse

Минимум:

- Snapshot PVC или `clickhouse-backup` в object storage.
- Проверка restore в отдельный ClickHouse pod.
- Так как query logs почти все отключены и TTL 1 день, не рассчитывайте на system logs как на audit trail.

### Redis

Redis сейчас пишет RDB `--save 60 1`. Делайте backup PVC или копируйте `/data/dump.rdb` в object storage. Если Redis хранит только cache, documented RPO может быть “можно потерять”. Если хранит очередь или state, нужен более строгий backup или Redis HA.

### NATS JetStream

Если JetStream хранит важные сообщения, включите persistent store на PVC и делайте snapshots. Для HA нужен NATS cluster. Одиночный NATS с PVC не защищает от отказа ноды.

### Etcd k3s

Настройте snapshots k3s etcd:

```bash
sudo mkdir -p /var/lib/rancher/k3s/server/db/snapshots
sudo k3s etcd-snapshot save --name manual-before-twir
sudo k3s etcd-snapshot ls
```

Для регулярных snapshots k3s имеет flags `--etcd-snapshot-schedule-cron`, `--etcd-snapshot-retention`, `--etcd-s3`. Лучше включить S3-compatible remote storage.

## Проверка после деплоя

### Cluster

```bash
kubectl get nodes -o wide
kubectl get nodes --show-labels
kubectl get pods -A
kubectl -n kube-system get pods
kubectl -n twir get all
```

Проверьте k3s/etcd:

```bash
sudo k3s kubectl get --raw='/readyz?verbose'
sudo k3s etcd-snapshot ls
```

### Storage

```bash
kubectl -n twir get pvc
kubectl -n twir describe pvc postgres-data
kubectl -n twir get pv -o wide
kubectl -n twir describe pod postgres-0
```

Убедитесь, что Postgres, ClickHouse и Temporal Postgres сидят на DB primary ноде, Redis и backup на cache/backup ноде.

### Secrets и ConfigMaps

```bash
kubectl -n twir get secrets
kubectl -n twir get configmaps
kubectl -n twir describe secret registry-twir-app
kubectl -n twir describe configmap clickhouse-config
```

Не выводите значения secrets в терминал, который логируется.

### Stateful dependencies

```bash
kubectl -n twir rollout status statefulset/postgres
kubectl -n twir rollout status statefulset/clickhouse
kubectl -n twir rollout status statefulset/redis
kubectl -n twir rollout status statefulset/temporal-postgres
kubectl -n twir logs statefulset/postgres --tail=100
kubectl -n twir logs statefulset/clickhouse --tail=100
```

Тест подключения:

```bash
kubectl -n twir run pg-client --rm -it --image=postgres:18.1 -- \
  psql -h postgres -U "$TWIR_POSTGRES_USER" -d "$TWIR_POSTGRES_DB" -c 'select 1;'

kubectl -n twir run redis-client --rm -it --image=redis:8.0.2 -- \
  redis-cli -h redis ping
```

### Migrations

```bash
kubectl -n twir apply -f migrations-job.yaml
kubectl -n twir wait --for=condition=complete job/migrations --timeout=10m
kubectl -n twir logs job/migrations
```

Если Job упал:

```bash
kubectl -n twir describe job migrations
kubectl -n twir logs job/migrations --previous
```

Не запускайте app services, пока причина падения migrations не исправлена.

### App rollout

```bash
kubectl -n twir rollout status deploy/api-gql
kubectl -n twir rollout status deploy/bots
kubectl -n twir rollout status deploy/parser
kubectl -n twir rollout status deploy/web
kubectl -n twir rollout status deploy/dashboard
kubectl -n twir rollout status deploy/overlays
kubectl -n twir rollout status deploy/websockets
kubectl -n twir get pods -o wide
kubectl -n twir get endpoints
```

Проверка логов:

```bash
kubectl -n twir logs deploy/api-gql --tail=100
kubectl -n twir logs deploy/websockets --tail=100
kubectl -n twir logs deploy/eventsub --tail=100
```

### DNS и routes

До cutover используйте `curl --resolve`:

```bash
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/dashboard/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/overlays/
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/api/health
curl -I --resolve twir.app:443:NODE1_IP https://twir.app/s/test
curl -I --resolve services-bots.twir.app:443:NODE1_IP https://services-bots.twir.app/
curl -I --resolve music-recognizer.twir.app:443:NODE1_IP https://music-recognizer.twir.app/
```

Проверьте custom-domain fallback:

```bash
curl -I --resolve some-custom-domain.example:443:NODE1_IP https://some-custom-domain.example/abc
```

Ожидаемый upstream path для fallback: `/v1/short-links/abc` в `api-gql`.

### WebSocket

Используйте `websocat` или `wscat`:

```bash
websocat -v wss://twir.app/socket
```

Проверьте, что соединение не рвется через обычный idle timeout ingress. Для nginx/Traefik выставьте read/send/idle timeout около `3600s`, как в HAProxy.

### Rate limit и compression

Rate limit:

```bash
for i in $(seq 1 1100); do curl -s -o /dev/null -w "%{http_code}\n" https://twir.app/; done | sort | uniq -c
```

Если используете не HAProxy, реализуйте аналогичный лимит через ingress controller, WAF или app/API gateway. Проверьте, что ключ берется из правильного real client IP header.

Compression:

```bash
curl -H 'Accept-Encoding: gzip' -I https://twir.app/
curl -H 'Accept-Encoding: gzip' -I https://twir.app/api/health
```

### Backup restore rehearsal

Перед production cutover выполните минимум один restore test:

```bash
kubectl -n twir get cronjob
kubectl -n twir create job --from=cronjob/postgres-backup postgres-backup-manual
kubectl -n twir logs job/postgres-backup-manual
```

Затем восстановите backup в отдельный временный Postgres и проверьте выборочные таблицы. Без restore rehearsal backup считается непроверенным.

## Traffic migration

1. Поднимите Kubernetes рядом со Swarm.
2. Синхронизируйте secrets и configs.
3. Восстановите свежую копию production DB или подключите Kubernetes apps к production DB только на короткий controlled test, если это безопасно.
4. Запустите migrations один раз, не параллельно со Swarm migrations.
5. Протестируйте Kubernetes через preview DNS или `curl --resolve`.
6. Уменьшите Swarm write activity или переведите сайт в maintenance, если схема БД меняется несовместимо.
7. Переключите DNS A/AAAA на Kubernetes ingress.
8. Следите за logs, 5xx, latency, Postgres connections, queue lag, NATS/Temporal health.
9. Держите Swarm готовым к rollback до стабилизации.

Если приложение не поддерживает active-active между Swarm и Kubernetes, не держите оба окружения одновременно пишущими в одну БД.

## Rollback plan

### До DNS cutover

Откат простой:

```bash
kubectl -n twir scale deploy --all --replicas=0
kubectl -n twir scale statefulset --all --replicas=0
```

Swarm остается production.

### После DNS cutover без миграции схемы

1. Верните DNS A/AAAA на старые Swarm IP.
2. Подождите TTL.
3. Остановите Kubernetes public ingress или app deployments, чтобы избежать split-brain.
4. Проверьте Swarm logs и routes.

### После миграции схемы

Rollback возможен только если есть обратимая миграция или backup restore plan.

1. Остановите write traffic.
2. Оцените, можно ли старый код работать с новой схемой.
3. Если нельзя, восстановите Postgres backup на Swarm или отдельный Postgres.
4. Верните DNS.
5. Зафиксируйте потерю данных по RPO, если restore откатывает последние записи.

Для каждого release храните:

- pinned image tags;
- примененные манифесты;
- номер DB migration;
- backup id перед миграцией;
- команду rollback для ingress и deployments.

Kubernetes rollout undo работает только если tag immutable:

```bash
kubectl -n twir rollout history deploy/api-gql
kubectl -n twir rollout undo deploy/api-gql --to-revision=2
```

С `latest` это ненадежно.

## Основные риски

- Базы на local PV не HA. Отказ DB ноды означает downtime и restore или ручной ремонт.
- Три ноды это минимум для etcd HA. При отказе одной ноды остается кворум 2 из 3, но capacity резко падает.
- Control plane и workloads делят ресурсы. Без requests/limits можно повредить etcd и базы нагрузкой приложений.
- `latest` tags ломают воспроизводимость deploy и rollback.
- `depends_on` из Swarm не переносится напрямую. Нужны probes, Jobs и deploy order.
- HAProxy features не все один к одному переносятся в стандартный Ingress. Особенно stick-table rate limit, custom fallback rewrite и cache.
- PgDog `users.toml` сейчас содержит статические примерные credentials. В Kubernetes это должен быть Secret.
- NATS JetStream в compose без volume. При переносе нужно явно решить, какие streams критичны и где хранится data dir.
- Adminer, Temporal UI, NATS metrics и HAProxy stats нельзя публиковать без auth или VPN.
- Cloudflare/real client IP headers должны совпасть с rate limit логикой, иначе лимит будет работать по IP proxy.
- Hetzner Cloud CSI применим только для Hetzner Cloud VMs в одной location. Для dedicated servers он не решает storage.

## Финальный checklist

Перед cutover:

- Все image tags pinned, нет `latest` для production Twir images.
- `kubectl get nodes` показывает 3 Ready ноды.
- Labels назначены и stateful pods сидят на нужных нодах.
- PVC Bound, storage path понятен и задокументирован.
- Secrets созданы, статические sample passwords удалены из runtime манифестов.
- ConfigMaps для ClickHouse, PgDog, OTel применены.
- Postgres, ClickHouse, Redis, Temporal Postgres, NATS готовы.
- PgDog и PgBouncer принимают подключения.
- Migrations Job завершился успешно.
- Все Deployments и StatefulSets rolled out.
- Ingress routes повторяют `/api`, `/s/`, `/dashboard`, `/overlays`, `/socket`, bots, music-recognizer и fallback.
- WebSocket держится дольше стандартного idle timeout.
- Rate limit и gzip проверены или явно заменены другой реализацией.
- Backup job успешно сделал backup.
- Restore rehearsal выполнен в отдельном окружении.
- Monitoring и alerts включены для nodes, pods, ingress, DB, Redis, NATS, backups.
- DNS TTL снижен.
- Rollback plan проверен и понятен дежурному.

После cutover:

- Проверить `twir.app`, `cf.twir.app`, `services-bots.twir.app`, `music-recognizer.twir.app`.
- Проверить login, dashboard, overlays, API, short links, WebSocket.
- Смотреть 5xx, latency, CPU/RAM, disk IO, Postgres connections, queue lag минимум 1 час.
- Не удалять Swarm окружение до успешного backup, restore rehearsal и стабильного production окна.
