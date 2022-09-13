package mynats

import (
	"time"

	"github.com/nats-io/nats.go"
)

func New(url string) (*nats.Conn, error) {
	return nats.Connect("nats://localhost:4222",
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.Name("Parser-go"),
		nats.Timeout(5*time.Second),
	)
}
