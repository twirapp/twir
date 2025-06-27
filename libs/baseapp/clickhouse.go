package baseapp

import (
	"context"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	config "github.com/satont/twir/libs/config"
)

type ClickhouseClient struct {
	driver.Conn
}

func NewClickHouse(appName string) func(config.Config) (*ClickhouseClient, error) {
	return func(cfg config.Config) (*ClickhouseClient, error) {
		options, err := clickhouse.ParseDSN(cfg.ClickhouseUrl)
		if err != nil {
			return nil, err
		}

		conn, err := clickhouse.Open(
			&clickhouse.Options{
				Addr: options.Addr,
				Auth: options.Auth,
				// TLS: &tls.Config{
				// 	InsecureSkipVerify: true,
				// },
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

		return &ClickhouseClient{conn}, nil
	}
}
