# Development

## Pre requirements

Before starting developing application you need these tools installed:

#### All system-wide dependencies provided by VSCode and Devcontainers

You can easy setup dependencies for project via installation of these deps:

- [Docker](https://docs.docker.com/engine/)
- [Visual Studio Code](https://code.visualstudio.com/)
- [VSCode Devcontainers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- Node 18
- pnpm
```bash
npm i -g pnpm
```

#### If you not using vscode

Oh, dear... You need to check `.devcontainer/Dockerfile` file and check what i'm installed via package manager
with versions from `docker-compose.dev.yml`. Also you need to check other tools installed in container (for example via
go install).

Sorry, i won't describe how to do that, because there is few deps, and they can be changed in anytime. I'm lazy update
readme, but Dockerfile must be always up-to-date.

Write command to run needed services
```bash
docker compose -f docker-compose.dev.yml up -d postgres redis tts
```
Notice, we omited `tsuwari` here, because it needed only for vscode.
### Next steps

Well, now we are almost ready for developing project, just few steps.

- Create twitch application https://dev.twitch.tv/console/apps
- Set `http://localhost:3005/login` and `https://twitchtokengenerator.com` as your redirect url's for twitch application
- Go to https://twitchtokengenerator.com, set clientID and clientSecret from your app and generate initial token WITH
  ALL SCOPES
- Read `.env`, I'm tried to describe important parts.
    ```bash
    cp .env.example .env
    ```
  Then fill that with values.
- Open project folder in devcontainer. Execute "Dev Containers: open folder in container" via vscode commands, or via
  another ways. Doesn't metter.
- Execute `pnpm install`

Now you are read to run project:

```bash
pnpm dev
```

And when everything starts open https://localhost:3005

#### Migrations

Migrations done via [goose](https://github.com/pressly/goose).
1. Navigate to folder
	```bash
	cd libs/migrations/migrations
	```
2. Use command for create new migration
	```bash
	goose create new_migration_name sql
	```
3. Run new created migrations (optional)
	```bash
	cd libs/migrations
	go run main.go
	```
##### Write `go` models

If you not familiar with the go, you can check existed models.

1. Go to `libs/gomodels`
2. Create new file and describe the go schema
3. Do not forget about `TableName()` for struct

##### How to run migrations

Migration ran automatically when you execute `pnpm dev`, or you can run them manually via
