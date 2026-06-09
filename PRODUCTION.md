# Deploy application to production

## Release model

- Production builds are triggered by git tags like `v1.0.0`.
- Docker images are published with the normalized tag, for example `v1.0.0` -> `1-0-0`.
- The webhook updates managed Swarm services directly through the Docker API.
- GitHub Actions finishes the build and then calls the deploy webhook running inside Swarm.

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
- `haproxy`

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
docker stack deploy -c docker-compose.stack.yml --with-registry-auth --prune twir
```

This first deploy also starts the `deploy-webhook` service.

After that, normal releases do not need repository files on the manager for deployment, because the webhook proxies image updates directly to the Docker Engine API over `/var/run/docker.sock`.

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

1. Build release images with tag `1-0-0`
2. Push them to `registry.twir.app`
3. Call the webhook
4. The webhook will update managed services in the `twir` swarm stack to image tag `1-0-0`

## Manual webhook call

```bash
curl -X POST \
  -H "Authorization: Bearer super-secret-webhook-token" \
  -H "Content-Type: application/json" \
  -d '{"service":"api-gql","imageTag":"1-0-0","refName":"v1.0.0","trigger":"manual"}' \
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
