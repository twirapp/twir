package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"
	cfg "tsuwari/config"
	"tsuwari/parser/internal/commands"
	"tsuwari/parser/internal/config/redis"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables"

	mynats "tsuwari/parser/internal/config/nats"

	twitch "tsuwari/parser/internal/config/twitch"
	natshandlers "tsuwari/parser/internal/handlers/nats"
	usersauth "tsuwari/parser/internal/twitch/user"

	"github.com/getsentry/sentry-go"
	"github.com/samber/lo"
	parserproto "github.com/satont/tsuwari/nats/parser"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"go.uber.org/zap"
	proto "google.golang.org/protobuf/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func main() {
	cfg, err := cfg.New()
	if err != nil || cfg == nil {
		fmt.Println(err)
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

	var logger *zap.Logger

	if cfg.AppEnv == "development" {
		l, _ := zap.NewDevelopment()
		logger = l
	} else {
		l, _ := zap.NewProduction()
		logger = l
	}

	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	r := redis.New(cfg.RedisUrl)
	defer r.Close()
	n, err := mynats.New(cfg.NatsUrl)
	defer n.Close()
	if err != nil {
		panic(err)
	}
	natsProtoConn, err := nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)
	defer natsProtoConn.Close()
	if err != nil {
		panic(err)
	}

	usersAuthService := usersauth.New(usersauth.UsersServiceOpts{
		Db:           db,
		ClientId:     cfg.TwitchClientId,
		ClientSecret: cfg.TwitchClientSecret,
	})
	twitchClient := twitch.New(*cfg)
	variablesService := variables.New()
	commandsService := commands.New(commands.CommandsOpts{
		Redis:            r,
		VariablesService: variablesService,
		Db:               db,
		UsersAuth:        usersAuthService,
		Nats:             n,
		Twitch:           twitchClient,
	})
	natsHandlers := natshandlers.New(natshandlers.NatsService{
		Redis:     r,
		Variables: variablesService,
		Commands:  commandsService,
	})

	if err != nil {
		panic(err)
	}

	natsProtoConn.QueueSubscribe("parser.handleProcessCommand", "parser", func(m *nats.Msg) {
		start := time.Now()
		data := parserproto.Request{}
		err := proto.Unmarshal(m.Data, &data)
		if err != nil {
			panic(err)
		}

		r := natsHandlers.HandleProcessCommand(data)

		if r != nil {
			res, _ := proto.Marshal(r)

			if err == nil {
				m.Respond(res)
			} else {
				fmt.Println(err)
			}
		} else {
			m.Respond([]byte{})
		}

		defer logger.Sugar().Infow("HandleProcessCommand ended.",
			"in", data.Message.Text,
			"out", r,
			"took", time.Since(start),
			"channelId", data.Channel.Id,
			"senderId", data.Sender.Id,
		)
		m.Ack()
	})

	natsProtoConn.QueueSubscribe("bots.getVariables", "parser", func(m *nats.Msg) {
		vars := lo.Map(variablesService.Store, func(v types.Variable, _ int) *parserproto.Variable {
			desc := v.Name
			if v.Description != nil {
				desc = *v.Description
			}
			example := v.Name
			if v.Example != nil {
				example = *v.Example
			}
			return &parserproto.Variable{
				Name:        v.Name,
				Example:     example,
				Description: desc,
			}
		})

		res, _ := proto.Marshal(&parserproto.GetVariablesResponse{
			List: vars,
		})

		m.Respond(res)
		m.Ack()
	})

	natsProtoConn.QueueSubscribe("bots.getDefaultCommands", "parser", func(m *nats.Msg) {
		list := make([]*parserproto.DefaultCommand, len(commandsService.DefaultCommands))

		for i, v := range commandsService.DefaultCommands {
			cmd := &parserproto.DefaultCommand{
				Name:        v.Name,
				Description: *v.Description,
				Visible:     v.Visible,
				Permission:  v.Permission,
				Module:      *v.Module,
				IsReply:     v.IsReply,
			}

			list[i] = cmd
		}

		res, _ := proto.Marshal(&parserproto.GetDefaultCommandsResponse{
			List: list,
		})

		m.Respond(res)
	})

	natsProtoConn.QueueSubscribe("parser.parseTextResponse", "parser", func(m *nats.Msg) {
		data := parserproto.ParseResponseRequest{}
		err := proto.Unmarshal(m.Data, &data)
		if err != nil {
			logger.Error(err.Error())
			return
		}
		text := natsHandlers.ParseResponse(data)
		bytes, err := proto.Marshal(&parserproto.ParseResponseResponse{
			Responses: []string{text},
		})
		if err != nil {
			logger.Error(err.Error())
			return
		}

		m.Respond(bytes)
		m.Ack()
	})

	logger.Info("Started")

	// runtime.Goexit()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Fatalf("Exiting")
}
