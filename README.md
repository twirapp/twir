# Development

## Requirements

- [Node.js (20+)](https://nodejs.org/en)
- [Pnpm](https://pnpm.io/)
- [Go (1.21+)](https://go.dev/)
- [Protobuf-compiler](https://grpc.io/docs/protoc-installation/)

> [!WARNING]
> Installation of protobuf depends on your system, google it.

- [Docker](https://docs.docker.com/engine/)
- Installed Go cli dependencies

> [!CAUTION]
> You need to setup [GOPATH](https://github.com/golang/go/wiki/SettingGOPATH) variable before executing lines below.

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/twitchtv/twirp/protoc-gen-twirp@latest
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Prepare

- Install Node.js dependencies
```bash
pnpm install --frozen-lockfile
```


- Run needed services (Postgres, Adminer, Redis, Minio)
```bash
docker compose -f docker-compose.dev.yml up -d
```

### Next steps

Well, now we are almost ready for developing project, just few steps.

- Create twitch application https://dev.twitch.tv/console/apps
- Set `http://localhost:3005/login` and `https://twitchtokengenerator.com` as your redirect url's for twitch application
- Go to https://twitchtokengenerator.com, set clientID and clientSecret from your app and generate initial token WITH
  ALL SCOPES
- `cp .env.example .env` and fill required envs
#### Now you are read to run project:

```bash
pnpm dev
```

And when everything starts open https://localhost:3005

## Writing migrations

Migrations done via [goose](https://github.com/pressly/goose).
1. Navigate to folder
	```bash
	cd libs/migrations/migrations
	```
2. Use command for create new migration
	```bash
	goose create new_migration_name sql
	```
* Run new created migrations (optional, because it's running when you execute `pnpm dev`)
	```bash
	cd libs/migrations
	go run main.go
	```
##### Write `go` models

1. Go to `libs/gomodels`
2. Create new file and describe the go schema
3. Do not forget about `TableName()` for struct
