# План рефакторинга apps/parser/ на PGX с идиоматичным DI

## Обзор текущего состояния

### Использование GORM в parser

GORM используется в `apps/parser` для следующих операций:

1. **cmd/main.go:155-169** - Инициализация GORM подключения к PostgreSQL
2. **internal/types/services/services.go:51** - GORM в структуре Services
3. **internal/services/chat_wall/chat_wall.go:187-205** - Запрос UsersStats для проверки
	 подписчиков/VIP
4. **internal/commands/permit/permit.go:67-83** - Создание permits через транзакцию
5. **internal/commands/games/** - Игровые команды (8ball, duel, russian roulette, seppuku)
6. **internal/commands/songrequest/youtube/sr.go** - Song request операции
7. **internal/commands/manage/** - Управление командами (add_aliase, edit_command, del_command)

### GORM модели используемые в parser

Основные модели из `libs/gomodels`:

- `UsersStats` - статистика пользователей
- `ChannelsPermits` - permits система
- `ChannelDuel` - дуэли
- `ChannelGames8Ball` - 8ball игра
- `ChannelGamesRussianRoulette` - русская рулетка
- `ChannelGamesSeppuku` - сеppуку игра
- `RequestedSong` - запрошенные песни
- `ChannelsCommands` - команды канала
- `ChannelsCommandsResponses` - ответы команд
- `ChannelsStreams` - стримы
- `ChannelsSongRequestsSettings` - настройки song requests

### Проблемы архитектуры

#### 1. GORM зависимости

- Тяжеловесный ORM со многими ненужными функциями
- Отсутствие compile-time проверки типов
- Необходимость перехода на pgx

#### 2. Команды (internal/commands/commands.go)

- Огромная функция `New()` с 60+ командами в одном месте
- Использует `map[string]*types.DefaultCommand` для хранения
- Большой метод `ParseCommandResponses` (500+ строк) со смешанной логикой
- Нет четкого разделения команд на группы/категории

#### 3. Переменные (internal/variables/variables.go)

- Огромная функция `New()` с 50+ переменными в одном месте
- Использует `map[string]*types.Variable` для хранения
- Нет четкого разделения переменных на группы/категории
- Сложная логика в `ParseVariablesInText`

#### 4. Кешер (internal/cacher/)

- Неправильное название: это Request-Scoped Cache, не просто "cacher"
- Использует ненужный интерфейс `types.DataCacher`
- Использует GORM в некоторых методах
- Разделен на 7 файлов для разных типов кешируемых данных
- Много индивидуальных мутексов для thread safety

## Цели миграции

1. **Удалить GORM** из `apps/parser`
2. **Использовать идиоматичные Go подходы для DI** (простые конструкторы без внешних библиотек)
3. **Использовать существующие pgx репозитории** где возможно
4. **Создать недостающие pgx репозитории** для отсутствующих таблиц
5. **Рефакторить commands.go** - разделить на логические группы
6. **Рефакторить variables.go** - разделить на логические группы
7. **Рефакторить cacher** - переименовать и упростить
8. **Обновить архитектуру** для соответствия существующим паттернам проекта

## План работы

### Этап 1: Создание недостающих pgx репозиториев

#### 1.1 Репозитории для систем модерации

- `libs/repositories/channels_permits/pgx/pgx.go`
- Интерфейс: `libs/repositories/channels_permits/repository.go`
- Модель: `libs/repositories/channels_permits/model/permits.go`

**Необходимые методы:**

- `Create(ctx, input)` - создать permit
- `GetMany(ctx, input)` - получить permits для канала
- `Delete(ctx, id)` - удалить permit
- `DeleteByUserID(ctx, channelID, userID)` - удалить permit пользователя

#### 1.2 Репозитории для игр

- `libs/repositories/channels_games_8ball/pgx/pgx.go`
- `libs/repositories/channels_games_russian_roulette/pgx/pgx.go`
- `libs/repositories/channels_games_seppuku/pgx/pgx.go`
- `libs/repositories/channels_games_duel/pgx/pgx.go`

**ChannelGames8Ball методы:**

- `GetByChannelID(ctx, channelID)` - получить настройки
- `Update(ctx, id, input)` - обновить настройки

**ChannelGamesRussianRoulette методы:**

- `GetByChannelID(ctx, channelID)` - получить настройки
- `Update(ctx, id, input)` - обновить настройки

**ChannelGamesSeppuku методы:**

- `GetByChannelID(ctx, channelID)` - получить настройки
- `Update(ctx, id, input)` - обновить настройки

**ChannelGamesDuel методы:**

- `Create(ctx, input)` - создать дуэль
- `GetByChannelID(ctx, channelID)` - получить настройки
- `GetActive(ctx, channelID)` - получить активную дуэль
- `Update(ctx, id, input)` - обновить дуэль
- `Delete(ctx, id)` - удалить дуэль

#### 1.3 Репозитории для song requests

- `libs/repositories/requested_songs/pgx/pgx.go`
- Интерфейс: `libs/repositories/requested_songs/repository.go`
- Модель: `libs/repositories/requested_songs/model/requested_song.go`

**Необходимые методы:**

- `Create(ctx, input)` - создать song request
- `GetLatestByChannelID(ctx, channelID)` - получить последний
- `GetManyByChannelID(ctx, channelID, input)` - получить очередь
- `GetByVideoID(ctx, channelID, videoID)` - проверить дубликат
- `CountByChannelID(ctx, channelID)` - получить количество
- `CountByUserID(ctx, channelID, userID)` - получить количество пользователя
- `Update(ctx, id, input)` - обновить (soft delete)
- `Delete(ctx, id)` - удалить

#### 1.4 Репозиторий для UsersStats

**ВНИМАНИЕ:** Уже существует `libs/repositories/userswithstats`, нужно проверить совместимость

Требуемые методы для parser:

- `GetManyByUserIDsAndChannelID(ctx, userIDs, channelID)` - получить статистику для нескольких
	пользователей

Если существующий репозиторий не покрывает эту функциональность, расширить его.

### Этап 2: Рефакторинг cacher → requestcache

#### 2.1 Переименование и реструктуризация

**Новое название:** `requestcache` (более точное - это request-scoped cache)

**Новая структура:**

```
apps/parser/internal/requestcache/
├── requestcache.go       # Основная структура и конструктор
├── cache.go             # Базовый интерфейс кеширования
├── integrations.go      # Кеш для интеграций
├── twitch.go           # Кеш для Twitch данных
├── faceit.go          # Кеш для Faceit данных
├── valorant.go        # Кеш для Valorant данных
├── song.go            # Кеш для текущей песни
└── subage.go          # Кеш для subage
```

#### 2.2 Удаление интерфейса types.DataCacher

**Действия:**

- Удалить файл `internal/types/variables_cacher.go`
- Удалить поле `Cacher DataCacher` из `types.ParseContext`
- Заменить на конкретные методы через requestcache

**Причина:** Интерфейс не нужен - мы используем один конкретный тип

#### 2.3 Удаление GORM из кешера

**Методы использующие GORM:**

1. `GetEnabledChannelIntegrations()` - использовать `channels_integrations.Repository`
2. `GetGbUserStats()` - использовать `userswithstats.Repository`
3. `GetTwitchUserFollow()` - использовать `channels.Repository`

#### 2.4 Упрощение структуры кеша

**Сейчас:**

```go
type locks struct {
stream      sync.Mutex
dbUserStats sync.Mutex
// ... 11 мутексов
}

type cache struct {
stream      *model.ChannelsStreams
dbUserStats *model.UsersStats
// ... 11 полей
}
```

**Предлагается:**

```go
type RequestCache struct {
mu    sync.RWMutex
data   map[cacheKey]interface{}
services *services.Services
channelID string
senderID   string
}

// Использовать map вместо отдельных полей
type cacheKey string

const (
cacheKeyStream cacheKey = "stream"
cacheKeyUserStats cacheKey = "user_stats"
cacheKeyTwitchUserFollow cacheKey = "twitch_user_follow"
// ... остальные ключи
)
```

**Преимущества:**

- Один мутекс вместо 11
- Единый API для всех типов данных
- Упрощение тестирования

#### 2.5 Создание нового requestcache.go

```go
package requestcache

import (
	"context"
	"sync"

	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/libs/integrations/seventv/api"
	"github.com/twirapp/twir/libs/twitch"
)

type RequestCache struct {
	mu        sync.RWMutex
	data      map[cacheKey]interface{}
	services  *services.Services
	channelID string
	senderID  string
}

type Opts struct {
	Services  *services.Services
	ChannelID string
	SenderID  string
}

func New(opts *Opts) *RequestCache {
	return &RequestCache{
		data:      make(map[cacheKey]interface{}),
		services:  opts.Services,
		channelID: opts.ChannelID,
		senderID:  opts.SenderID,
	}
}

// Методы кеширования для каждого типа данных
```

### Этап 3: Рефакторинг commands.go

#### 3.1 Разделение команд на логические группы

**Текущая структура:**

```go
func New(opts *Opts) *Commands {
commands := lo.SliceToMap(
[]*types.DefaultCommand{
song.CurrentSong,
song.Playlist,
// ... 60+ команд в одном месте
}, ...
)
}
```

**Новая структура:**

```
apps/parser/internal/commands/
├── commands.go            # Основная структура и Registry
├── registry.go           # Registry для регистрации команд
├── commands/             # Определения команд
│   ├── moderation/       # !permit, !nuke, !shoutout
│   ├── music/            # !sr, !srlist, !skip, !playlist
│   ├── games/            # !8ball, !duel, !roulette, !seppuku
│   ├── stats/            # !stats, !uptime, !userage, !watchtime
│   ├── tts/              # !tts say, !tts skip, !tts voices
│   ├── overlays/         # !brb, !kappagen
│   ├── clips/            # !clip
│   └── markers/         # !marker
└── handler.go            # Общая логика обработки команд
```

#### 3.2 Создание Registry паттерна

```go
package commands

import "github.com/twirapp/twir/apps/parser/internal/types"

type CommandRegistry struct {
	commands map[string]*types.DefaultCommand
}

func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*types.DefaultCommand),
	}
}

func (r *CommandRegistry) Register(cmd *types.DefaultCommand) {
	r.commands[cmd.Name] = cmd
}

func (r *CommandRegistry) Get(name string) (*types.DefaultCommand, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *CommandRegistry) GetAll() map[string]*types.DefaultCommand {
	return r.commands
}
```

#### 3.3 Разделение регистрации по группам

```go
package commands

import (
	moderationpkg "github.com/twirapp/twir/apps/parser/internal/commands/moderation"
	musicpkg "github.com/twirapp/twir/apps/parser/internal/commands/music"
	gamespkg "github.com/twirapp/twir/apps/parser/internal/commands/games"
	// ... остальные группы
)

func RegisterAllCommands(registry *CommandRegistry) {
	// Регистрация по группам
	moderationpkg.Register(registry)
	musicpkg.Register(registry)
	gamespkg.Register(registry)
	// ...
}
```

#### 3.4 Упрощение Commands структуры

```go
type Commands struct {
registry         *CommandRegistry
services         *services.Services
variablesService *variables.Variables
}

type Opts struct {
Services         *services.Services
VariablesService *variables.Variables
}

func New(opts *Opts) *Commands {
registry := NewRegistry()
RegisterAllCommands(registry)

return &Commands{
registry:         registry,
services:         opts.Services,
variablesService: opts.VariablesService,
}
}
```

#### 3.5 Рефакторинг ParseCommandResponses

**Проблема:** Метод содержит 500+ строк с смешанной логикой

**Решение:** Разделить на отдельные методы:

```go
func (c *Commands) ParseCommandResponses(
ctx context.Context,
command *FindByMessageResult,
requestData twitch.TwitchChatMessage,
userRoles []model.ChannelRole,
userChannelStats *model.UsersStats,
dbUser *model.Users,
) *busparser.CommandParseResponse {
// Основной метод - координация
}

func (c *Commands) parseDefaultCommand(...) *busparser.CommandParseResponse
func (c *Commands) parseCustomCommand(...) *busparser.CommandParseResponse
func (c *Commands) filterResponsesByCategory(...) []commandresponsemodel.Response
func (c *Commands) replaceVariablesInResponses(...) []string
```

### Этап 4: Рефакторинг variables.go

#### 4.1 Разделение переменных на логические группы

**Текущая структура:**

```go
func New(opts *Opts) *Variables {
store := lo.SliceToMap(
[]*types.Variable{
command_param.Variable,
commands_list.Variable,
// ... 50+ переменных в одном месте
}, ...
)
}
```

**Новая структура:**

```
apps/parser/internal/variables/
├── variables.go          # Основная структура и Registry
├── registry.go          # Registry для регистрации переменных
├── variables/          # Определения переменных
│   ├── user/           # $(user), $(user.age), $(user.messages)
│   ├── stream/         # $(stream), $(stream.title), $(stream.uptime)
│   ├── command/        # $(command), $(command.count)
│   ├── song/           # $(song), $(song.history)
│   ├── emotes/         # $(emotes), $(emotes.top)
│   ├── external/        # $(faceit), $(valorant), $(7tv)
│   └── utility/        # $(random.number), $(repeat), $(request)
└── parser.go            # Парсинг переменных из текста
```

#### 4.2 Создание Registry для переменных

```go
package variables

import "github.com/twirapp/twir/apps/parser/internal/types"

type VariableRegistry struct {
	variables map[string]*types.Variable
}

func NewRegistry() *VariableRegistry {
	return &VariableRegistry{
		variables: make(map[string]*types.Variable),
	}
}

func (r *VariableRegistry) Register(v *types.Variable) {
	r.variables[v.Name] = v
}

func (r *VariableRegistry) Get(name string) (*types.Variable, bool) {
	v, ok := r.variables[name]
	return v, ok
}

func (r *VariableRegistry) GetByPriority() []*types.Variable {
	sorted := make([]*types.Variable, 0, len(r.variables))
	for _, v := range r.variables {
		sorted = append(sorted, v)
	}
	sort.Slice(
		sorted, func(i, j int) bool {
			return sorted[j].Priority < sorted[i].Priority
		}
	)
	return sorted
}
```

#### 4.3 Разделение регистрации по группам

```go
package variables

import (
	userpkg "github.com/twirapp/twir/apps/parser/internal/variables/user"
	streampkg "github.com/twirapp/twir/apps/parser/internal/variables/stream"
	commandpkg "github.com/twirapp/twir/apps/parser/internal/variables/command"
	// ... остальные группы
)

func RegisterAllVariables(registry *VariableRegistry) {
	userpkg.Register(registry)
	streampkg.Register(registry)
	commandpkg.Register(registry)
	// ...
}
```

#### 4.4 Упрощение Variables структуры

```go
type Variables struct {
registry   *VariableRegistry
services   *services.Services
}

type Opts struct {
Services *services.Services
}

func New(opts *Opts) *Variables {
registry := NewRegistry()
RegisterAllVariables(registry)

return &Variables{
registry: registry,
services: opts.Services,
}
}
```

#### 4.5 Рефакторинг ParseVariablesInText

**Проблема:** Метод содержит сложную логику с goroutines и WaitGroup

**Решение:** Упростить и разделить:

```go
func (v *Variables) ParseVariablesInText(
ctx context.Context,
parseCtx *types.ParseContext,
input string,
) []string {
variables := v.registry.GetByPriority()
return v.processVariables(ctx, parseCtx, input, variables)
}

func (v *Variables) processVariables(
ctx context.Context,
parseCtx *types.ParseContext,
input string,
variables []*types.Variable,
) []string {
// Логика обработки переменных
}
```

### Этап 5: Рефакторинг Services и DI

#### 5.1 Удаление GORM из Services

Обновить `internal/types/services/services.go`:

- Удалить поле `Gorm *gorm.DB`
- Добавить новые репозитории:
	- `ChannelsPermitsRepo channels_permits.Repository`
	- `Games8BallRepo channels_games_8ball.Repository`
	- `GamesRussianRouletteRepo channels_games_russian_roulette.Repository`
	- `GamesSeppukuRepo channels_games_seppuku.Repository`
	- `GamesDuelRepo channels_games_duel.Repository`
	- `RequestedSongsRepo requested_songs.Repository`
	- `UsersWithStatsRepo userswithstats.Repository` (расширить если нужно)

#### 5.2 Использование идиоматичного DI (Go конструкторы)

**Без внешних библиотек:**

```go
// cmd/main.go
func initServices(
ctx context.Context,
cfg *cfg.Config,
pgxPool *pgxpool.Pool,
redisClient *redis.Client,
bus *buscore.Bus,
) (*services.Services, error) {
// Создание репозиториев
repositories, err := initRepositories(pgxPool)
if err != nil {
return nil, err
}

// Создание кешей
caches := initCaches(repositories, bus, redisClient)

// Создание gRPC клиентов
grpcClients := initGrpcClients(cfg)

// Создание сервисов
return initServices(cfg, repositories, caches, grpcClients, bus, redisClient), nil
}

func initRepositories(pgxPool *pgxpool.Pool) (*Repositories, error) {
return &Repositories{
ChannelsPermitsRepo:     channels_permits.NewFx(pgxPool),
Games8BallRepo:         channels_games_8ball.NewFx(pgxPool),
GamesRussianRouletteRepo: channels_games_russian_roulette.NewFx(pgxPool),
GamesSeppukuRepo:       channels_games_seppuku.NewFx(pgxPool),
GamesDuelRepo:          channels_games_duel.NewFx(pgxPool),
RequestedSongsRepo:     requested_songs.NewFx(pgxPool),
// ... остальные репозитории
}, nil
}

func initCaches(repositories *Repositories, bus *buscore.Bus, redisClient *redis.Client) *Caches {
return &Caches{
CommandsCache:      commandscache.New(repositories.CommandsRepo, bus),
CommandsPrefixCache: commandsprefixcache.New(repositories.CommandsPrefixRepo, bus),
// ... остальные кеши
}
}

func initServices(cfg *cfg.Config, repos *Repositories, caches *Caches, grpc *GrpcClients, bus *buscore.Bus, redis *redis.Client) *services.Services {
return &services.Services{
Config:     cfg,
// ... репозитории
// ... кеши
GrpcClients: grpc,
Bus:        bus,
Redis:      redisClient,
// ... остальные поля
}
}
```

#### 5.3 Обновление инициализации в main.go

**Удалить:**

- `gorm.io/driver/postgres` импорт
- `gorm.io/gorm` импорт
- Инициализацию GORM (строки 148-169)
- `sqlDb := stdlib.OpenDBFromPool(pgxconn)`

**Заменить на:**

```go
pgxconn, err := pgxpool.NewWithConfig(context.Background(), connConfig)
if err != nil {
panic(err)
}
defer pgxconn.Close()

services, err := initServices(ctx, config, pgxconn, redisClient, bus)
if err != nil {
panic(err)
}
```

### Этап 6: Обновление команд (Commands)

#### 6.1 Обновить internal/commands/permit/permit.go

- Заменить `parseCtx.Services.Gorm.WithContext(ctx).Transaction(...)` на использование репозитория
- Использовать `channels_permits.Repository`

#### 6.2 Обновить internal/commands/games/

- `8ball.go` - использовать `channels_games_8ball.Repository`
- `russian_roulette.go` - использовать `channels_games_russian_roulette.Repository`
- `seppuku.go` - использовать `channels_games_seppuku.Repository`
- `duel.go`, `duel_accept.go`, `duel_handler.go` - использовать `channels_games_duel.Repository`

#### 6.3 Обновить internal/commands/songrequest/youtube/

- `sr.go` - использовать `requested_songs.Repository`
- `skip.go` - использовать `requested_songs.Repository`

#### 6.4 Обновить internal/commands/manage/

- `add_aliase.go`, `edit_command.go`, `del_command.go` - использовать существующий
	`commands_repository`

#### 6.5 Обновить internal/commands/dota/

Проверить закомментированный GORM код и либо удалить его, либо обновить для нужной функциональности.

### Этап 7: Обновление Chat Wall Service

#### 7.1 Обновить internal/services/chat_wall/chat_wall.go

- Удалить поле `Gorm *gorm.DB`
- Заменить запрос `UsersStats` (строки 187-205) на:
	```go
	usersStats, err := c.usersWithStatsRepo.GetManyByUserIDsAndChannelID(ctx, userIDs, input.ChannelID)
	```

#### 7.2 Обновить инициализацию Chat Wall Service

- Удалить `Gorm: db` из Opts
- Добавить `UsersWithStatsRepo userswithstats.Repository`

### Этап 8: Обновление переменных (Variables)

#### 8.1 Проверить internal/variables/donations/

- `last_donate/last_donate.go`
- `top_donate/top_donate.go`
- `top_donate_stream/top_donate_stream.go`

Если они используют GORM, обновить для использования pgx репозиториев.

#### 8.2 Обновить использование types.DataCacher

Заменить все использования `parseCtx.Cacher.GetXxx()` на соответствующие методы из requestcache

### Этап 9: Очистка зависимостей

#### 9.1 Удалить из go.mod

```bash
go mod tidy
```

Проверить что `gorm.io/gorm` и `gorm.io/driver/postgres` больше не используются.

#### 9.2 Удалить libs/gomodels (частично)

Проверить все места использования `libs/gomodels` в parser и мигрировать на pgx репозитории.

**Примечание:** Не удалять `libs/gomodels` полностью, так как он может использоваться в других
частях проекта (apps/api-gql, etc.).

#### 2.2 Создание DI файла

Создать `apps/parser/internal/di/providers.go`:

```go
//go:generate go tool kessoku $GOFILE
package di

import (
	"github.com/mazrean/kessoku"
)

var _ = kessoku.Inject[*services.Services](
	"InitServices",
	// Async providers для независимых сервисов
	kessoku.Async(kessoku.Provide(NewPostgresPool)),
	kessoku.Async(kessoku.Provide(NewRedisClient)),
	kessoku.Async(kessoku.Provide(NewRepositories)),
	kessoku.Async(kessoku.Provide(NewCaches)),
	kessoku.Async(kessoku.Provide(NewGrpcClients)),
	kessoku.Async(kessoku.Provide(NewBusClients)),

	// Синхронные провайдеры для зависимых сервисов
	kessoku.Provide(NewServices),
)
```

#### 2.3 Создание провайдеров

- `NewPostgresPool(*cfg.Config)` - создает pgxpool
- `NewRedisClient(*cfg.Config)` - создает redis client
- `NewRepositories(pgxpool.Pool, *redis.Client)` - создает все репозитории
- `NewCaches(repositories, bus, redis)` - создает кеши
- `NewGrpcClients(*cfg.Config)` - создает gRPC клиенты
- `NewBusClients(*cfg.Config)` - создает bus клиенты
- `NewServices(...)` - создает Services struct

### Этап 3: Рефакторинг Services

#### 3.1 Удаление GORM из Services

Обновить `internal/types/services/services.go`:

- Удалить поле `Gorm *gorm.DB`
- Добавить новые репозитории:
	- `ChannelsPermitsRepo channels_permits.Repository`
	- `Games8BallRepo channels_games_8ball.Repository`
	- `GamesRussianRouletteRepo channels_games_russian_roulette.Repository`
	- `GamesSeppukuRepo channels_games_seppuku.Repository`
	- `GamesDuelRepo channels_games_duel.Repository`
	- `RequestedSongsRepo requested_songs.Repository`
	- `UsersWithStatsRepo userswithstats.Repository` (расширить если нужно)

#### 3.2 Обновление инициализации в main.go

Удалить:

- `gorm.io/driver/postgres` импорт
- `gorm.io/gorm` импорт
- Инициализацию GORM (строки 148-169)
- `sqlDb := stdlib.OpenDBFromPool(pgxconn)`

Использовать Kessoku:

```go
services, err := di.InitServices(ctx, config)
if err != nil {
panic(err)
}
```

### Этап 4: Обновление команд (Commands)

#### 4.1 Обновить internal/commands/permit/permit.go

- Заменить `parseCtx.Services.Gorm.WithContext(ctx).Transaction(...)` на использование репозитория
- Использовать `channels_permits.Repository`

#### 4.2 Обновить internal/commands/games/

- `8ball.go` - использовать `channels_games_8ball.Repository`
- `russian_roulette.go` - использовать `channels_games_russian_roulette.Repository`
- `seppuku.go` - использовать `channels_games_seppuku.Repository`
- `duel.go`, `duel_accept.go`, `duel_handler.go` - использовать `channels_games_duel.Repository`

#### 4.3 Обновить internal/commands/songrequest/youtube/

- `sr.go` - использовать `requested_songs.Repository`
- `skip.go` - использовать `requested_songs.Repository`

#### 4.4 Обновить internal/commands/manage/

- `add_aliase.go`, `edit_command.go`, `del_command.go` - использовать существующий
	`commands_repository`

#### 4.5 Обновить internal/commands/dota/

Проверить закомментированный GORM код и либо удалить его, либо обновить для нужной функциональности.

### Этап 5: Обновить Chat Wall Service

#### 5.1 Обновить internal/services/chat_wall/chat_wall.go

- Удалить поле `Gorm *gorm.DB`
- Заменить запрос `UsersStats` (строки 187-205) на:
	```go
	usersStats, err := c.usersWithStatsRepo.GetManyByUserIDsAndChannelID(ctx, userIDs, input.ChannelID)
	```

#### 5.2 Обновить инициализацию Chat Wall Service

- Удалить `Gorm: db` из Opts
- Добавить `UsersWithStatsRepo userswithstats.Repository`

### Этап 6: Обновить переменные (Variables)

#### 6.1 Проверить internal/variables/donations/

- `last_donate/last_donate.go`
- `top_donate/top_donate.go`
- `top_donate_stream/top_donate_stream.go`

Если они используют GORM, обновить для использования pgx репозиториев.

### Этап 7: Очистка зависимостей

#### 7.1 Удалить из go.mod

```bash
go mod tidy
```

Проверить что `gorm.io/gorm` и `gorm.io/driver/postgres` больше не используются.

#### 7.2 Удалить libs/gomodels (частично)

Проверить все места использования `libs/gomodels` и мигрировать на pgx репозитории.

**Примечание:** Не удалять `libs/gomodels` полностью, так как он может использоваться в других
частях проекта (apps/api-gql, etc.).

### Этап 8: Тестирование

#### 8.1 Unit тесты

- Написать тесты для новых pgx репозиториев
- Написать тесты для обновленных команд

#### 8.2 Интеграционные тесты

- Протестировать permit систему
- Протестировать игровые команды
- Протестировать song requests
- Протестировать chat wall

#### 8.3 Ручное тестирование

- Запустить parser локально
- Протестировать все обновленные команды

## Детали реализации pgx репозиториев

### Паттерн (согласно существующему коду)

```go
package pgx

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/twirapp/twir/libs/repositories/your_repo"
	"github.com/twirapp/twir/libs/repositories/your_repo/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

var _ your_repo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

// Методы репозитория с использованием pgx и squirrel
```

### Создание entities

Согласно `AGENTS.md`:

- Создать entity в `libs/entities/{entity_name}/entity.go`
- Entity должен содержать только domain логику и validation
- Использовать Nil паттерн:

```go
type YourEntity struct {
ID        string
ChannelID string
// другие поля
isNil bool
}

func (e YourEntity) IsNil() bool {
return e.isNil
}

var Nil = &YourEntity{
isNil: true,
}
```

## Детали реализации requestcache

### Базовый паттерн кеширования

```go
package requestcache

import (
	"context"
	"sync"
)

type Cache interface {
	Get(ctx context.Context, key interface{}) (interface{}, error)
	Set(ctx context.Context, key, value interface{}) error
}

// Дефолтная реализация с sync.Map
type InMemoryCache struct {
	data sync.Map
}

func (c *InMemoryCache) Get(_ context.Context, key interface{}) (interface{}, error) {
	value, ok := c.data.Load(key)
	if !ok {
		return nil, ErrCacheNotFound
	}
	return value, nil
}

func (c *InMemoryCache) Set(_ context.Context, key, value interface{}) error {
	c.data.Store(key, value)
	return nil
}
```

### Пример кеширования Twitch User

```go
type TwitchUserCache struct {
cache Cache
}

func (t *TwitchUserCache) GetByID(ctx context.Context, userID string) (*helix.User, error) {
cached, err := t.cache.Get(ctx, "twitch_user_by_id:"+userID)
if err == nil {
return cached.(*helix.User), nil
}

// Fetch from API
user, err := t.twitchClient.GetUsers(&helix.UsersParams{
IDs: []string{userID},
})

// Cache result
_ = t.cache.Set(ctx, "twitch_user_by_id:"+userID, user)

return user, err
}
```

## Детали реализации Registry паттерна

### Command Registry

```go
package commands

type CommandRegistry struct {
	commands map[string]*types.DefaultCommand
}

func NewRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]*types.DefaultCommand),
	}
}

func (r *CommandRegistry) Register(cmd *types.DefaultCommand) {
	r.commands[cmd.Name] = cmd
}

func (r *CommandRegistry) RegisterGroup(cmds []*types.DefaultCommand) {
	for _, cmd := range cmds {
		r.Register(cmd)
	}
}

func (r *CommandRegistry) Find(name string) (*types.DefaultCommand, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}
```

### Variable Registry

```go
package variables

type VariableRegistry struct {
	variables map[string]*types.Variable
}

func NewRegistry() *VariableRegistry {
	return &VariableRegistry{
		variables: make(map[string]*types.Variable),
	}
}

func (r *VariableRegistry) Register(v *types.Variable) {
	r.variables[v.Name] = v
}

func (r *VariableRegistry) RegisterGroup(vars []*types.Variable) {
	for _, v := range vars {
		r.Register(v)
	}
}

func (r *VariableRegistry) Get(name string) (*types.Variable, bool) {
	v, ok := r.variables[name]
	return v, ok
}

func (r *VariableRegistry) GetByPriority() []*types.Variable {
	sorted := make([]*types.Variable, 0, len(r.variables))
	for _, v := range r.variables {
		sorted = append(sorted, v)
	}
	sort.Slice(
		sorted, func(i, j int) bool {
			return sorted[j].Priority < sorted[i].Priority
		}
	)
	return sorted
}
```

## Требуемые миграции базы данных

### Создать миграции для новых таблиц

```bash
bun cli m create --name create_channels_permits_table --db postgres --type sql
bun cli m create --name create_channels_games_duel_table --db postgres --type sql
bun cli m create --name create_requested_songs_table --db postgres --type sql
bun cli m create --name create_channels_games_8ball_table --db postgres --type sql
bun cli m create --name create_channels_games_russian_roulette_table --db postgres --type sql
bun cli m create --name create_channels_games_seppuku_table --db postgres --type sql
```

**ВНИМАНИЕ:** Проверить существующие таблицы перед созданием миграций!

## Преимущества миграции

1. **Чистота кода:** Удаление GORM упрощает зависимости и код
2. **Типобезопасность:** pgx предоставляет compile-time типобезопасность
3. **Контроль:** Полный контроль над SQL запросами через squirrel
4. **Консистентность:** Единая архитектура с остальной частью проекта
5. **Меньше зависимостей:** Удаление тяжеловесного GORM
6. **Улучшенный код:** Registry паттерн улучшает читаемость commands/variables
7. **Лучшая структура:** Группировка команд/переменных упрощает навигацию
8. **Request-scoped cache:** Упрощение и улучшение системы кеширования
9. **Идиоматичный Go:** Использование конструкторов вместо внешних DI фреймворков

## Порядок выполнения (по приоритету)

1. ✅ Анализ текущего состояния (сделан)
2. ⬜ Создать entity для новых таблиц
3. ⬜ Создать pgx репозитории для отсутствующих таблиц
4. ⬜ Рефакторинг requestcache (переименование, реструктуризация)
5. ⬜ Рефакторинг commands.go (Registry паттерн, разделение на группы)
6. ⬜ Рефакторинг variables.go (Registry паттерн, разделение на группы)
7. ⬜ Обновить Services struct (удаление GORM, добавление репозиториев)
8. ⬜ Обновить Chat Wall Service
9. ⬜ Обновить commands/permit
10. ⬜ Обновить commands/games
11. ⬜ Обновить commands/songrequest
12. ⬜ Обновить commands/manage
13. ⬜ Обновить variables/donations
14. ⬜ Реализовать идиоматичный DI (конструкторы)
15. ⬜ Обновить main.go
16. ⬜ Удалить GORM зависимости
17. ⬜ Тестирование
18. ⬜ Очистка кода

## Дополнительные заметки

### Взаимодействие с существующими репозиториями

- `commands_repository` уже существует и может использоваться для commands
- `channels_info_history_repo` уже существует
- `channel_event_lists_repo` уже существует
- `channels_games_voteban_repo` уже существует

### Команды с комментариями GORM

Некоторые команды в `internal/commands/dota/` имеют закомментированный код с GORM. Этот код нужно
либо:

1. Удалить полностью если не используется
2. Обновить для использования pgx если нужен функционал

### UsersStats существующий репозиторий

Проверить `libs/repositories/userswithstats` - он может уже покрывать нужную функциональность для
parser. Если нет - расширить его.

### Идиоматичный DI в Go

Согласно стандартным Go практикам:

- **Использовать простые функции-конструкторы**
- **Не использовать кодогенерацию для DI** (как Kessoku)
- **Передавать зависимости явно через параметры функций**
- **Создавать зависимые компоненты по цепочке**

Пример:

```go
func NewServices(
repos *Repositories,
caches *Caches,
grpc *GrpcClients,
bus *buscore.Bus,
redis *redis.Client,
cfg *config.Config,
) *Services {
return &Services{
Repos: repos,
Caches: caches,
GrpcClients: grpc,
Bus: bus,
Redis: redis,
Config: cfg,
}
}
```

## Вывод

Этот план позволяет систематически мигрировать `apps/parser` с GORM на чистую архитектуру с pgx и
идиоматичным Go DI, сохраняя функциональность и улучшая качество кода через рефакторинг commands.go,
variables.go и cacher системы.
