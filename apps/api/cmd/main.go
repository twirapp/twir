package main

import (
	"log"
	"time"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/locales/en_US"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/satont/go-helix/v2"
	auth "github.com/satont/tsuwari/apps/api/internal/api/auth"
	apiv1 "github.com/satont/tsuwari/apps/api/internal/api/v1"
	"github.com/satont/tsuwari/apps/api/internal/middlewares"
	"github.com/satont/tsuwari/apps/api/internal/types"
	myNats "github.com/satont/tsuwari/libs/nats"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	cfg "github.com/satont/tsuwari/libs/config"

	"github.com/satont/tsuwari/apps/api/internal/services/redis"

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
	d, _ := db.DB()
	d.SetMaxOpenConns(20)
	d.SetConnMaxIdleTime(1 * time.Minute)

	natsEncodedConn, natsConn, err := myNats.New(cfg.NatsUrl)
	if err != nil {
		panic(err)
	}
	defer natsEncodedConn.Close()

	storage := redis.NewCache(cfg.RedisUrl)

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

	/* app.Use(func(c *fiber.Ctx) error {
		c.Next()
		defer logger.Sugar().Infow("incoming request",
			"method", c.Method(),
			"path", c.Path(),
			"code", c.Response().StatusCode(),
		)
		return nil
	}) */
	appLogger, _ := zap.NewDevelopment()
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: appLogger,
	}))

	v1 := app.Group("/v1")

	services := types.Services{
		DB:                  db,
		RedisStorage:        storage,
		Validator:           validator,
		ValidatorTranslator: transEN,
		Twitch: twitch.NewClient(&helix.Options{
			ClientID:     cfg.TwitchClientId,
			ClientSecret: cfg.TwitchClientSecret,
			RedirectURI:  cfg.TwitchCallbackUrl,
		}),
		Logger: logger,
		Cfg:    cfg,
		Nats:   natsConn,
	}

	if cfg.FeedbackTelegramBotToken != nil {
		services.TgBotApi, _ = tgbotapi.NewBotAPI(*cfg.FeedbackTelegramBotToken)
	}

	apiv1.Setup(v1, services)
	auth.Setup(app, services)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not found")
	})

	log.Fatal(app.Listen(":3002"))
}
