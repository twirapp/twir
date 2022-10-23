package main

import (
	"log"
	"runtime"
	"tsuwari/twitch"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/storage/redis"
	apiv1 "github.com/satont/tsuwari/apps/api-go/internal/api/v1"
	"github.com/satont/tsuwari/apps/api-go/internal/middlewares"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	myNats "github.com/satont/tsuwari/libs/nats"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cfg "tsuwari/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	logger, _ := zap.NewDevelopment()
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		logger.Sugar().Error(err)
		panic("Cannot load config of application")
	}

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

	natsEncodedConn, natsConn, err := myNats.New(cfg.NatsUrl)
	if err != nil {
		panic(err)
	}
	defer natsEncodedConn.Close()

	store := redis.New(redis.Config{
		URL:       cfg.RedisUrl,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	validator := validator.New()
	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	enTranslations.RegisterDefaultTranslations(validator, transEN)
	errorMiddleware := middlewares.ErrorHandler(transEN, logger)

	app := fiber.New(fiber.Config{
		ErrorHandler: errorMiddleware,
	})
	app.Use(compress.New())
	app.Use(fiberlogger.New())

	v1 := app.Group("/v1")

	services := types.Services{
		DB:                  db,
		RedisStorage:        store,
		Validator:           validator,
		ValidatorTranslator: transEN,
		Twitch:              twitch.NewClient(cfg.TwitchClientId, cfg.TwitchClientSecret),
		Logger:              logger,
		Cfg:                 cfg,
		Nats:                natsConn,
	}

	if cfg.FeedbackTelegramBotToken != nil {
		services.TgBotApi, _ = tgbotapi.NewBotAPI(*cfg.FeedbackTelegramBotToken)
	}

	apiv1.Setup(v1, services)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not found")
	})

	log.Fatal(app.Listen(":3002"))
}
