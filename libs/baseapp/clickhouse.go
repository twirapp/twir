package baseapp

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	config "github.com/satont/twir/libs/config"

	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
)

func NewClickHouse(appName string) func(config.Config) (*twirclickhouse.ClickhouseClient, error) {
	return func(cfg config.Config) (*twirclickhouse.ClickhouseClient, error) {
		options, err := clickhouse.ParseDSN(cfg.ClickhouseUrl)
		if err != nil {
			return nil, err
		}

		conn, err := clickhouse.Open(
			&clickhouse.Options{
				Addr: options.Addr,
				Auth: options.Auth,
				Settings: clickhouse.Settings{
					"max_execution_time":    30,
					"async_insert":          "1",
					"wait_for_async_insert": "1",
				},
				DialTimeout: time.Second * 5,
				Compression: &clickhouse.Compression{
					Method: clickhouse.CompressionLZ4,
				},
				Debug:                false,
				BlockBufferSize:      10,
				MaxCompressionBuffer: 10240,
				ClientInfo: clickhouse.ClientInfo{
					Products: []struct {
						Name    string
						Version string
					}{
						{Name: appName},
					},
				},
			},
		)

		pingCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := conn.Ping(pingCtx); err != nil {
			return nil, err
		}

		return twirclickhouse.New(conn), nil
	}
}
