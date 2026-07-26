package messagehandler

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	chatwallrepository "github.com/twirapp/twir/libs/repositories/chat_wall"
	chatwallmodel "github.com/twirapp/twir/libs/repositories/chat_wall/model"
)

type messageIDCache struct{}

func (messageIDCache) Get(context.Context, string) kv.Valuer { return messageIDMissingValue{} }

func (messageIDCache) Set(context.Context, string, any, ...kvoptions.Option) error { return nil }

func (messageIDCache) SetMany(context.Context, []kv.SetMany) error { return nil }

func (messageIDCache) Delete(context.Context, string) error { return nil }

func (messageIDCache) DeleteMany(context.Context, []string) error { return nil }

func (messageIDCache) Exists(context.Context, string) (bool, error) { return false, nil }

func (messageIDCache) ExistsMany(context.Context, []string) ([]bool, error) { return nil, nil }

func (messageIDCache) GetKeysByPattern(context.Context, string) ([]string, error) { return nil, nil }

type messageIDMissingValue struct{}

func (messageIDMissingValue) Int() (int64, error) { return 0, kv.ErrKeyNil }

func (messageIDMissingValue) String() (string, error) { return "", kv.ErrKeyNil }

func (messageIDMissingValue) Bytes() ([]byte, error) { return nil, kv.ErrKeyNil }

func (messageIDMissingValue) Bool() (bool, error) { return false, kv.ErrKeyNil }

func (messageIDMissingValue) Float() (float64, error) { return 0, kv.ErrKeyNil }

func (messageIDMissingValue) Scan(any) error { return kv.ErrKeyNil }

func (messageIDMissingValue) Err() error { return kv.ErrKeyNil }

type messageIDRedisRecorder struct {
	mu          sync.Mutex
	members     map[string]bool
	expirations map[string]bool
}

func newMessageIDRedisClient(t interface{ Cleanup(func()) }) (*redis.Client, *messageIDRedisRecorder) {
	recorder := &messageIDRedisRecorder{
		members:     make(map[string]bool),
		expirations: make(map[string]bool),
	}
	client := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	client.AddHook(recorder)
	t.Cleanup(func() { _ = client.Close() })
	return client, recorder
}

func (r *messageIDRedisRecorder) DialHook(next redis.DialHook) redis.DialHook { return next }

func (r *messageIDRedisRecorder) ProcessHook(redis.ProcessHook) redis.ProcessHook {
	return func(_ context.Context, command redis.Cmder) error {
		r.process(command)
		return nil
	}
}

func (r *messageIDRedisRecorder) ProcessPipelineHook(redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(_ context.Context, commands []redis.Cmder) error {
		for _, command := range commands {
			r.process(command)
		}
		return nil
	}
}

func (r *messageIDRedisRecorder) process(command redis.Cmder) {
	args := command.Args()
	if len(args) < 3 {
		return
	}

	key := fmt.Sprint(args[1])
	member := fmt.Sprint(args[2])

	r.mu.Lock()
	defer r.mu.Unlock()

	switch command.Name() {
	case "sismember":
		command.(*redis.BoolCmd).SetVal(r.members[key+"\x00"+member])
	case "sadd":
		for _, value := range args[2:] {
			r.members[key+"\x00"+fmt.Sprint(value)] = true
		}
	case "hset":
		r.members[key+"\x00"+member] = true
	case "hexpire":
		r.expirations[key+"\x00"+fmt.Sprint(args[len(args)-1])] = true
	}
}

type messageIDChatWallRepository struct{}

func (messageIDChatWallRepository) GetChannelSettings(context.Context, uuid.UUID) (chatwallmodel.ChatWallSettings, error) {
	return chatwallmodel.ChatWallSettings{}, nil
}

func (messageIDChatWallRepository) UpdateChannelSettings(context.Context, chatwallrepository.UpdateChannelSettingsInput) error {
	return nil
}

func (messageIDChatWallRepository) GetByID(context.Context, uuid.UUID) (chatwallmodel.ChatWall, error) {
	return chatwallmodel.ChatWall{}, nil
}

func (messageIDChatWallRepository) GetMany(context.Context, chatwallrepository.GetManyInput) ([]chatwallmodel.ChatWall, error) {
	return nil, nil
}

func (messageIDChatWallRepository) GetLogs(context.Context, uuid.UUID) ([]chatwallmodel.ChatWallLog, error) {
	return nil, nil
}

func (messageIDChatWallRepository) Create(context.Context, chatwallrepository.CreateInput) (chatwallmodel.ChatWall, error) {
	return chatwallmodel.ChatWall{}, nil
}

func (messageIDChatWallRepository) CreateLog(context.Context, chatwallrepository.CreateLogInput) error {
	return nil
}

func (messageIDChatWallRepository) CreateManyLogs(context.Context, []chatwallrepository.CreateLogInput) error {
	return nil
}

func (messageIDChatWallRepository) Update(
	context.Context,
	uuid.UUID,
	chatwallrepository.UpdateInput,
) (chatwallmodel.ChatWall, error) {
	return chatwallmodel.ChatWall{}, nil
}

func (messageIDChatWallRepository) Delete(context.Context, uuid.UUID) error { return nil }
