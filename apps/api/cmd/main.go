package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/swagger"
	"github.com/satont/tsuwari/apps/api/internal/api/auth"
	apiv1 "github.com/satont/tsuwari/apps/api/internal/api/v1"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/clients"

	"github.com/satont/tsuwari/apps/api/internal/services/redis"

	"github.com/satont/tsuwari/apps/api/internal/services/timers_service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/satont/tsuwari/apps/api/docs"
	gormLogger "gorm.io/gorm/logger"
)

// @title Tsuwari api
// @version 1.0
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

	zap.ReplaceGlobals(logger)

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

	dbConnOpts, err := pq.ParseURL(cfg.DatabaseUrl)
	if err != nil {
		panic(fmt.Errorf("cannot parse postgres url connection: %w", err))
	}
	sqlxConn, err := sqlx.Connect("postgres", dbConnOpts)
	if err != nil {
		log.Fatalln(err)
	}

	validator := validator.New()
	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	enTranslations.RegisterDefaultTranslations(validator, transEN)
	errorMiddleware := middlewares.ErrorHandler(logger.Sugar(), transEN)
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

	v1 := app.Group("/v1")

	neededServices := &types.Services{
		Logger:              logger.Sugar(),
		Redis:               redis.New(cfg.RedisUrl),
		Gorm:                db,
		Sqlx:                sqlxConn,
		Config:              &config.Config{},
		RedisStorage:        redis.NewCache(cfg.RedisUrl),
		Validator:           validator,
		ValidatorTranslator: transEN,
		TgBotApi:            &tgbotapi.BotAPI{},
		Grpc: &types.GrpcClientsService{
			Integrations: clients.NewIntegrations(cfg.AppEnv),
			Parser:       clients.NewParser(cfg.AppEnv),
			EventSub:     clients.NewEventSub(cfg.AppEnv),
			Scheduler:    clients.NewScheduler(cfg.AppEnv),
			Timers:       clients.NewTimers(cfg.AppEnv),
			Bots:         clients.NewBots(cfg.AppEnv),
			Tokens:       clients.NewTokens(cfg.AppEnv),
		},
		TimersService: timers_service.NewTimersService(db, logger.Sugar()),
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
	logger.Sugar().Info("Api exiting...")
}
