# Development

## Requirements

* [Bun (1.2.2+)](https://bun.sh)
* [Go (1.21+)](https://go.dev/)

* [Docker](https://docs.docker.com/engine/)

## Cli

> [!NOTE]
> For MOST of project management tasks we use own written cli. You can use `bun cli help` for print cli usage

* Run needed services (Postgres, Adminer, Redis, Minio)
```bash
docker compose -f docker-compose.dev.yml up -d
```

* Install dependencies
```bash
bun cli deps
```

* Build libs
```bash
bun run build libs
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
bun dev
```
* Visit https://localhost:3005

## Writing migrations

* Use command for create new migration
```bash
bun cli migrations create
```
* Navigate to folder and edit new migration file
```bash
cd libs/migrations/migrations
```

* Run new created migrations (optional, because it's running when you execute `pnpm dev`)
```bash
bun cli migrations run
```
##### Write `go` models

* Go to `libs/gomodels`
* Create new file and describe the go schema
* Do not forget about `TableName()` for struct

## Https on localhost (optional)

We'll use `twir.localhost` domain, which is enables ability to grant ssl out of the box, but you can use any other domain and deal with ssl yourself.

* Add `https://twir.localhost/login` to your twitch application redirect url's

* Edit `.env`, change site base url and protocol for twitch callback:
```ini
SITE_BASE_URL=twir.localhost
TWITCH_CALLBACKURL=https://$SITE_BASE_URL/login
```

* Start application as usual:
```bash
bun dev
```

* Open https://twir.localhost
