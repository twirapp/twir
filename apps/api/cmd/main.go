package main

import (
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"

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
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	config "github.com/satont/tsuwari/libs/config"

	"github.com/satont/tsuwari/apps/api/internal/services/redis_storage"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/satont/tsuwari/apps/api/docs"
	_ "github.com/satont/tsuwari/apps/api/internal/di"
	gormLogger "gorm.io/gorm/logger"
)

// @title Fiber Example API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3002
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name api-key
// @description "apiKey" from /v1/profile response
func main() {
	logger, _ := zap.NewDevelopment()
	cfg, err := config.New()
	if err != nil || cfg == nil {
		logger.Sugar().Error(err)
		panic("Cannot load config of application")
	}

	do.ProvideValue[interfaces.Logger](di.Injector, logger.Sugar())
	do.ProvideValue(di.Injector, cfg)

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

	do.ProvideValue(di.Injector, db)
	do.ProvideValue(di.Injector, redis_storage.NewCache(cfg.RedisUrl))

	validator := validator.New()
	do.ProvideValue(di.Injector, validator)

	en := en_US.New()
	uni := ut.New(en, en)
	transEN, _ := uni.GetTranslator("en_US")
	do.ProvideValue(di.Injector, transEN)

	enTranslations.RegisterDefaultTranslations(validator, transEN)
	errorMiddleware := middlewares.ErrorHandler(transEN, logger)
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

	// appLogger, _ := zap.NewDevelopment()
	// app.Use(fiberzap.New(fiberzap.Config{
	// 	Logger: appLogger,

	// }))

	do.ProvideValue(di.Injector, clients.NewBots(cfg.AppEnv))
	do.ProvideValue(di.Injector, clients.NewTimers(cfg.AppEnv))
	do.ProvideValue(di.Injector, clients.NewScheduler(cfg.AppEnv))
	do.ProvideValue(di.Injector, clients.NewParser(cfg.AppEnv))
	do.ProvideValue(di.Injector, clients.NewEventSub(cfg.AppEnv))
	do.ProvideValue(di.Injector, clients.NewIntegrations(cfg.AppEnv))

	v1 := app.Group("/v1")

	do.ProvideValue(di.Injector, twitch.NewClient(&helix.Options{
		ClientID:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
		RedirectURI:  cfg.TwitchCallbackUrl,
	}))

	if cfg.FeedbackTelegramBotToken != nil {
		bot, _ := tgbotapi.NewBotAPI(*cfg.FeedbackTelegramBotToken)
		do.ProvideValue(di.Injector, bot)
	}

	apiv1.Setup(v1)
	auth.Setup(app)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).SendString("Not found")
	})

	go app.Listen(":3002")

	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	fmt.Println("Closing...")
	di.Injector.Shutdown()

	d, _ = db.DB()
	d.Close()
}
