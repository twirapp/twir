package seventv

import (
	"context"

	"github.com/kr/pretty"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm *gorm.DB
}

func New(opts Opts) error {
	s := Service{
		sockets: nil,
		gorm:    opts.Gorm,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.start(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.stop()
			},
		},
	)

	return nil
}

type Service struct {
	sockets []*wsConnection

	gorm *gorm.DB
}

func (c *Service) start(ctx context.Context) error {
	var channels []model.Channel
	if err := c.gorm.Select("id").Where(`"isEnabled" = true`).Find(&channels).Error; err != nil {
		return err
	}

	for _, channel := range channels {
		var conn *wsConnection
		for _, ws := range c.sockets {
			if len(ws.channels) < ws.channelsLimit {
				conn = ws
			}
		}
		if conn == nil {
			newConn, err := createConn(ctx, c.onMessage)
			if err != nil {
				return err
			}
			conn = newConn
			c.sockets = append(c.sockets, newConn)
		}

		conn.addChannel(ctx, channel.ID)
	}

	return nil
}

func (c *Service) stop() error {
	for _, socket := range c.sockets {
		if err := socket.Close(); err != nil {
			return err
		}
	}
	c.sockets = nil

	return nil
}

func (c *Service) onMessage(msg []byte) {
	pretty.Println(string(msg))
}
