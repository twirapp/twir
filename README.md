# Development

## Requirements

* [Node.js (20+)](https://nodejs.org/en)
* [Pnpm](https://pnpm.io/)
* [Go (1.21+)](https://go.dev/)

* [Docker](https://docs.docker.com/engine/)

## Cli

> [!NOTE]
> For MOST of project management tasks we use own written cli. You can use `pnpm cli help` for print cli usage

* Run needed services (Postgres, Adminer, Redis, Minio)
```bash
docker compose -f docker-compose.dev.yml up -d
```

* Install dependencies
```bash
pnpm cli deps
```

* Build libs
```
pnpm cli build libs
```

### Configure project for development

Well, now we are almost ready for developing project, just few steps.

* Create twitch application https://dev.twitch.tv/console/apps
* Set `http://localhost:3005/login` and `https://tokens-generator.twir.app` as your redirect url's for twitch application
* Go to https://tokens-generator.twir.app, set clientID and clientSecret from your app and generate initial token WITH
  ALL SCOPES
* `cp .env.example .env` and fill required envs

### Run project

* Start dev mode
```bash
pnpm cli dev
```
* Visit https://localhost:3005

## Writing migrations

* Use command for create new migration
```bash
pnpm cli migrations create
```
* Navigate to folder and edit new migration file
```bash
cd libs/migrations/migrations
```

* Run new created migrations (optional, because it's running when you execute `pnpm dev`)
```bash
pnpm cli migrations run
```
##### Write `go` models

* Go to `libs/gomodels`
* Create new file and describe the go schema
* Do not forget about `TableName()` for struct

## Https on localhost (optional)

We'll use `dev.twir.app` domain, but you can use any other domain.

* Add `https://twir.localhost/login` to your twitch application redirect url's

* Edit `.env` entries:
```ini
TWITCH_CALLBACKURL=https://twir.localhost/login
SITE_BASE_URL=twir.localhost
USE_WSS=true
```

* Start caddy:
```bash
pnpm cli proxy
```

* Open https://twir.localhost
