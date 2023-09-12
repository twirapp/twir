# Deploy application to production

### Setup manager

```bash
docker swarm init
```

### Setup node

```bash
docker swarm join --token token --advertise-addr nodepi  managerip:2377
```

### Create stack file

Just copy `docker-compose.stack.yml` to manager filesystem.

### Create network

```bash
docker network create -d overlay --attachable twir
```

### Create secrets

1. Doppler

   ```bash
   echo "token" | docker secret create twir_doppler_token -
   ```

2. Postgres

   ```bash
   echo "twir" | docker secret create twir_postgres_user -
   echo "twir" | docker secret create twir_postgres_db -
   echo "twir" | docker secret create twir_postgres_password -
   ```

### Deploy

This command will deploy services from stack to the cluster.

```bash
docker stack deploy -c docker-compose.yml --with-registry-auth --prune twir
```

### Update

Same command as deploy
