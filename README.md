# Development

## Requirements

* [Bun (1.2.19+)](https://bun.sh)
* [Go (1.24+)](https://go.dev/)

* [Docker](https://docs.docker.com/engine/)


### Development

> [!NOTE]
> For MOST of project management tasks we use own written cli. You can use `bun cli help` for print cli usage

* Create twitch application https://dev.twitch.tv/console/apps
* Set `http://localhost:3005/login` and `https://tokens-generator.twir.app` as your redirect url's for twitch application
* Go to https://tokens-generator.twir.app, set clientID and clientSecret from your app and generate initial token WITH
  ALL SCOPES
* `cp .env.example .env` and fill required envs

* Run needed services (Postgres, Adminer, Redis, Minio, e.t.c)
```bash
docker compose -f docker-compose.dev.yml up -d
```

* Start project
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

* Run new created migrations (optional, because it's running when you execute `bun dev`)
```bash
bun cli migrations run
```

## Https on localhost (optional)

We'll use `twir.localhost` domain, which is enables ability to grant ssl out of the box, but you can use any other domain and deal with ssl yourself.

* Add `https://twir.localhost/login` to your twitch application redirect url's

* Edit `.env`, change site base url:
```ini
SITE_BASE_URL=https://twir.localhost
```

* Start application as usual:
```bash
bun dev
```

* Open https://twir.localhost
