package main

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/services"
	"github.com/satont/tsuwari/libs/grpc/generated/bots"
	"github.com/satont/tsuwari/libs/grpc/generated/eventsub"
	"github.com/satont/tsuwari/libs/grpc/generated/integrations"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"github.com/satont/tsuwari/libs/grpc/generated/scheduler"
	"github.com/satont/tsuwari/libs/grpc/generated/timers"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/services"

	"github.com/gofiber/swagger"
	"github.com/satont/tsuwari/libs/grpc/clients"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/api/auth"
	apiv1 "github.com/satont/tsuwari/apps/api/internal/api/v1"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/satont/tsuwari/libs/config"

	"github.com/satont/tsuwari/apps/api/internal/services/redis"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/satont/tsuwari/apps/api/docs"
	gormLogger "gorm.io/gorm/logger"
)

// @title Tsuwari api
// @description Non-public api for tsuwari
// @host localhost:3002
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-key
func main() {
	logger, _ := zap.NewDevelopment()
	cfg, err := config.New()
	if err != nil || cfg == nil {
		logger.Sugar().Error(err)
		panic("Cannot load config of application")
	}

	do.ProvideValue[interfaces.Logger](di.Injector, logger.Sugar())

	if cfg.SentryDsn != "" {
		sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.SentryDsn,
			Environment:      cfg.AppEnv,
			Debug:            true,
			TracesSampleRate: 1.0,
		})
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		logger.Sugar().Error(err)
		panic("failed to connect database")
	}
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	do.ProvideValue[*gorm.DB](di.Injector, db)
	do.ProvideValue[interfaces.TimersService](di.Injector, services.NewTimersService())
	do.ProvideValue[*config.Config](di.Injector, cfg)

	storage := redis.NewCache(cfg.RedisUrl)

	validator := validator.New()
	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	enTranslations.RegisterDefaultTranslations(validator, transEN)
	errorMiddleware := middlewares.ErrorHandler(transEN)
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	app := fiber.New(fiber.Config{
		ErrorHandler: errorMiddleware,
	})
	app.Use(cors.New())
	if cfg.AppEnv == "development" {
		app.Get("/swagger/*", swagger.New(swagger.Config{
			URL:                  "http://localhost:3002/swagger/doc.json",
			DeepLinking:          false,
			DocExpansion:         "list",
			PersistAuthorization: true,
			Title:                "Tsuwari api",
			TryItOutEnabled:      true,
		}))
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	app.Use(compress.New())

	do.ProvideValue[integrations.IntegrationsClient](di.Injector, clients.NewIntegrations(cfg.AppEnv))
	do.ProvideValue[parser.ParserClient](di.Injector, clients.NewParser(cfg.AppEnv))
	do.ProvideValue[eventsub.EventSubClient](di.Injector, clients.NewEventSub(cfg.AppEnv))
	do.ProvideValue[scheduler.SchedulerClient](di.Injector, clients.NewScheduler(cfg.AppEnv))
	do.ProvideValue[timers.TimersClient](di.Injector, clients.NewTimers(cfg.AppEnv))
	do.ProvideValue[bots.BotsClient](di.Injector, clients.NewBots(cfg.AppEnv))

	v1 := app.Group("/v1")

	neededServices := types.Services{
		DB:                  db,
		RedisStorage:        storage,
		Validator:           validator,
		ValidatorTranslator: transEN,
		Twitch: twitch.NewClient(&helix.Options{
			ClientID:     cfg.TwitchClientId,
			ClientSecret: cfg.TwitchClientSecret,
			RedirectURI:  cfg.TwitchCallbackUrl,
		}),
		Cfg: cfg,
	}

	if cfg.FeedbackTelegramBotToken != nil {
		neededServices.TgBotApi, _ = tgbotapi.NewBotAPI(*cfg.FeedbackTelegramBotToken)
	}

	apiv1.Setup(v1, neededServices)
	auth.Setup(app, neededServices)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not found")
	})

	go app.Listen(":3002")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")

	d, _ = db.DB()
	d.Close()
}
