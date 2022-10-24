package myNats

import (
	"time"

	Nats "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
)

func New(natsUrl string) (*Nats.EncodedConn, *Nats.Conn, error) {
	n, err := Nats.Connect(natsUrl,
		Nats.RetryOnFailedConnect(true),
		Nats.MaxReconnects(10),
		Nats.ReconnectWait(time.Second),
		Nats.Name("Parser-go"),
		Nats.Timeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, err
	}

	natsProtoConn, err := Nats.NewEncodedConn(n, protobuf.PROTOBUF_ENCODER)
	if err != nil {
		return nil, nil, err
	}

	return natsProtoConn, n, nil
}
