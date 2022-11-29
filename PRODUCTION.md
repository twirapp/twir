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
docker network create -d overlay --attachable tsuwari
```

### Create secrets

1. Doppler

   ```bash
   echo "token" | docker secret create tsuwari_doppler_token -
   ```

2. Postgres

   ```bash
   echo "tsuwari" | docker secret create tsuwari_postgres_user -
   echo "tsuwari" | docker secret create tsuwari_postgres_db -
   echo "tsuwari" | docker secret create tsuwari_postgres_password -
   ```

### Deploy

This command will deploy services from stack to the cluster.

```bash
docker stack deploy -c docker-compose.yml --with-registry-auth --prune tsuwari
```

### Update

Same command as deploy
