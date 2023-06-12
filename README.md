# Development

## Prerequirements

Before starting developing application you need thos tools installed:

#### All system wide dependencies provided by VSCode and Devcontainers

You can easy setup dependencies for project via installation of this deps:

- [Docker](https://docs.docker.com/engine/)
- [Visual Studio Code](https://code.visualstudio.com/)
- [VSCode Devcontainers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
- Node 18
- pnpm
```bash
npm i -g pnpm
```

#### If you not using vscode

Oh, dear... You need to check `.devcontainer/Dockerfile` file and check what i'm installed via `pacman` package manager
with versions from `docker-compose.dev.yml`. Also you need to check other tools installed in container (for example via
go install).

Sorry, i won't describe how to do that, because there is few deps, and they can be changed in anytime. I'm lazy update
readme, but Dockerfile must be always up-to-date.

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

Swagger available at http://localhost:3005/api/swagger or http://localhost:3002/swagger

#### Adding new entities, models (migrations)

Migrations done via `typeorm`. So at first you always need to change `nodejs` entities, generate migrations, then
describe `go` models

##### Write `nodejs` models

If you not familar with nodejs, you can check existed entities.

1. Describe entity into `libs/typeorm/src/entities` folder. Also there is some example how i'm doing that
2. Add entity classname to `libs/typeorm/src/index.ts` into `entities` array. Thats how `typeorm` working in ESM mode
3. For generate migrations go to the typeorm folder `cd libs/typeorm` and
   run `pnpm migration:generate -n NameForMigration`

##### Write `go` models

If you not familar with the go, you can check existed models.

1. Go to `libs/gomodels`
2. Create new file and describe the go schema
3. Do not forget about `TableName()` for struct

##### How to run migrations

Migration runned automatically when you execute `task dev`, or you can run them manually via

```bash
cd libs/typeorm
pnpm run deploy
```
