# Deploy application to production

## Release model

- Production builds are triggered by git tags like `v1.0.0`.
- Docker images are published with the normalized tag, for example `v1.0.0` -> `1-0-0`.
- The webhook updates managed Swarm services directly through the Docker API.
- GitHub Actions finishes the build and then calls the deploy webhook running inside Swarm.
- CI builds and deploys **only affected** applications based on changed files (see "Affected builds" below).

## Setup manager

```bash
docker swarm init
```

## Setup node

```bash
docker swarm join --token token --advertise-addr nodepi managerip:2377
```

## Create networks

```bash
docker network create -d overlay --attachable twir
```

Create the other external networks used by the stack if they do not already exist:

- `traefik-public`
- `cloudflared`

## Create secrets

### App secrets

```bash
echo "token" | docker secret create twir_doppler_token -
echo "twir" | docker secret create twir_postgres_user -
echo "twir" | docker secret create twir_postgres_db -
echo "twir" | docker secret create twir_postgres_password -
```

### Deploy webhook secret

```bash
echo "super-secret-webhook-token" | docker secret create twir_deploy_webhook_token -
```

## First deploy

For the first deploy, choose the image tag that already exists in the registry and export it before deploying.

```bash
export TWIR_IMAGE_TAG=1-0-0
docker stack deploy --resolve-image changed -c docker-compose.stack.yml --with-registry-auth --prune twir
```

This first deploy also starts the `deploy-webhook` service.

After that, normal releases do not need repository files on the manager for deployment, because the webhook proxies image updates directly to the Docker Engine API over `/var/run/docker.sock`.

## Updating the stack configuration

When you change `docker-compose.stack.yml` (add/remove services, change labels, environment, replicas, etc.) **without changing image tags**, use `--resolve-image changed` to avoid unnecessary image pulls:

```bash
docker stack deploy --resolve-image changed -c docker-compose.stack.yml --with-registry-auth --prune twir
```

`--resolve-image changed` tells Docker to only check the registry for services whose image reference changed in the compose file. This means:
- Label/environment/replica changes → no image pulls → fast deploy
- New service added → image will be pulled → image must exist in registry
- Image tag changed in compose → image will be pulled

## Affected builds

CI determines which applications are affected by changed files using dependency graph analysis:

- Changes in `libs/config/` → rebuilds all Go apps that transitively depend on `config`
- Changes in `libs/frontend-chat/` → rebuilds `dashboard` and `overlays`
- Changes in `apps/api-gql/` → rebuilds only `api-gql`
- Changes in root files (`go.work`, `package.json`, `bun.lock`) → rebuilds all apps

The affected list is also passed to the deploy webhook, so only affected services are updated in the swarm.

To manually check which apps would be affected:

```bash
bun cli affected --files "libs/config/entity.go,apps/api-gql/main.go" --output json
```

## Two deploy mechanisms

| Mechanism | When to use | What it does |
|---|---|---|
| **Deploy webhook** | After CI builds new images | Updates specific services with new image tag via Docker API |
| **docker stack deploy** | Compose file changes (labels, env, new services) | Applies stack configuration; use `--resolve-image changed` to skip image pulls |

Normal workflow:
1. Push a tag → CI builds affected images → webhook updates affected services
2. Change compose file → run `docker stack deploy --resolve-image changed` manually

## GitHub configuration

Add these GitHub Actions secrets:

- `DOCKER_TWIR_LOGIN`
- `DOCKER_TWIR_PASSWORD`
- `DEPLOY_WEBHOOK_URL` for example `https://deploy.example.com/deploy` or `http://manager-ip:8090/deploy`
- `DEPLOY_WEBHOOK_TOKEN` matching `twir_deploy_webhook_token`

## Release

Create and push a git tag:

```bash
git tag v1.0.0
git push origin v1.0.0
```

The workflow will:

1. Determine which apps are affected by changes since the previous tag
2. Build only affected images with tag `1-0-0` (+ `latest`)
3. Push them to `registry.twir.app`
4. Call the webhook with the list of affected services
5. The webhook will update only affected services in the `twir` swarm stack to image tag `1-0-0`

## Manual webhook call

Deploy specific services:

```bash
curl -X POST \
  -H "Authorization: Bearer super-secret-webhook-token" \
  -H "Content-Type: application/json" \
  -d '{"imageTag":"1-0-0","refName":"v1.0.0","trigger":"manual","services":["api-gql","bots"]}' \
  http://manager-ip:8090/deploy
```

Deploy a single service:

```bash
curl -X POST \
  -H "Authorization: Bearer super-secret-webhook-token" \
  -H "Content-Type: application/json" \
  -d '{"service":"api-gql","imageTag":"1-0-0","refName":"v1.0.0","trigger":"manual"}' \
  http://manager-ip:8090/deploy
```

Deploy all services (omit `service`/`services`):

```bash
curl -X POST \
  -H "Authorization: Bearer super-secret-webhook-token" \
  -H "Content-Type: application/json" \
  -d '{"imageTag":"1-0-0","refName":"v1.0.0","trigger":"manual"}' \
  http://manager-ip:8090/deploy
```

## Health and status

```bash
curl http://manager-ip:8090/healthz
curl http://manager-ip:8090/status
```

## Manual deploy by tag

Inside the webhook container or from the same image:

```bash
docker exec -it $(docker ps -q -f name=deploy-webhook) /app/twir-cli deploy-apply --service api-gql --image-tag 1-0-0
```

# Postgres tunel

ssh -fN -L 54322:postgres:5432 satont@ssh.twir.app
