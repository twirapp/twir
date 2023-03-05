package pubsub

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type PubSub struct {
	ctx        context.Context
	publisher  *redis.Client
	subscriber *redis.Client
	subs       map[string]CB
}

func NewPubSub(url string) (*PubSub, error) {
	redisUrl, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	publisher := redis.NewClient(redisUrl)
	subscriber := redis.NewClient(redisUrl)

	return &PubSub{
		ctx:        context.Background(),
		publisher:  publisher,
		subscriber: subscriber,
	}, nil
}

type CB func(data []byte)

func (c *PubSub) Publish(topic string, data []byte) {
	c.publisher.Publish(c.ctx, topic, data)
}

func (c *PubSub) Subscribe(topic string, cb CB) {
	subscription := c.subscriber.Subscribe(c.ctx, topic)

	go func() {
		for {
			msg, err := subscription.ReceiveMessage(c.ctx)
			if err != nil {
				panic(err)
			}

			cb([]byte(msg.Payload))
		}
	}()
}
