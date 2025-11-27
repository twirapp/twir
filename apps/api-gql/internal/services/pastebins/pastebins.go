package pastebins

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goccy/go-json"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/pastebins"
	pastebinsmodel "github.com/twirapp/twir/libs/repositories/pastebins/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repo   pastebins.Repository
	KV     kv.KV
	Logger *slog.Logger
}

func New(opts Opts) *Service {
	return &Service{
		repo:   opts.Repo,
		kv:     opts.KV,
		logger: opts.Logger,
	}
}

type Service struct {
	repo   pastebins.Repository
	kv     kv.KV
	logger *slog.Logger
}

var ErrNotFound = fmt.Errorf("pastebin not found")

func (c *Service) mapToEntity(m pastebinsmodel.Pastebin) entity.Pastebin {
	return entity.Pastebin{
		ID:          m.ID,
		CreatedAt:   m.CreatedAt,
		Content:     m.Content,
		ExpireAt:    m.ExpireAt,
		OwnerUserID: m.OwnerUserID,
	}
}

func makeKvStoreKey(id string) string {
	return "twir:cache:pastebins:" + id
}

func (c *Service) deleteIfNeed(ctx context.Context, p entity.Pastebin) (bool, error) {
	if p.ExpireAt != nil && p.ExpireAt.Before(time.Now()) {
		return true, c.Delete(ctx, p.ID)
	}

	return false, nil
}

func (c *Service) GetByID(ctx context.Context, id string) (entity.Pastebin, error) {
	cacheKey := makeKvStoreKey(id)
	cachedBytes, err := c.kv.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, kv.ErrKeyNil) {
		return entity.PastebinNil, err
	}
	if len(cachedBytes) > 0 {
		var bin entity.Pastebin
		if err := json.Unmarshal(cachedBytes, &bin); err != nil {
			return entity.PastebinNil, err
		}

		if deleted, err := c.deleteIfNeed(ctx, bin); err != nil {
			return entity.PastebinNil, err
		} else if deleted {
			return entity.PastebinNil, ErrNotFound
		}

		return bin, nil
	}

	bin, err := c.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pastebins.ErrNotFound) {
			return entity.PastebinNil, ErrNotFound
		}

		return entity.PastebinNil, err
	}

	if deleted, err := c.deleteIfNeed(ctx, c.mapToEntity(bin)); err != nil {
		return entity.PastebinNil, err
	} else if deleted {
		return entity.PastebinNil, ErrNotFound
	}

	var cacheTime time.Duration
	if bin.ExpireAt != nil {
		cacheTime = bin.ExpireAt.Sub(time.Now())
	} else {
		cacheTime = 1 * 24 * time.Hour
	}

	converted := c.mapToEntity(bin)

	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		bytes, err := json.Marshal(converted)
		if err != nil {
			c.logger.Error("cannot convert pastebin entity to bytes", logger.Error(err))
			return
		}
		if err := c.kv.Set(cacheCtx, cacheKey, bytes, kvoptions.WithExpire(cacheTime)); err != nil {
			c.logger.Error("cannot save pastebin entity to kv", logger.Error(err))
			return
		}
	}()

	return converted, nil
}

type CreateInput struct {
	Content     string
	ExpireAt    *time.Time
	OwnerUserID *string
}

func (c *Service) generateID() string {
	return gonanoid.Must(5)
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Pastebin, error) {
	bin, err := c.repo.Create(
		ctx,
		pastebins.CreateInput{
			ID:          c.generateID(),
			Content:     input.Content,
			ExpireAt:    input.ExpireAt,
			OwnerUserID: input.OwnerUserID,
		},
	)
	if err != nil {
		return entity.PastebinNil, err
	}

	return c.mapToEntity(bin), nil
}

type GetManyInput struct {
	Page        int
	PerPage     int
	OwnerUserID string
}

type GetManyOutput struct {
	Items []entity.Pastebin
	Total int
}

func (c *Service) GetUserMany(ctx context.Context, input GetManyInput) (
	GetManyOutput,
	error,
) {
	bins, err := c.repo.GetManyByOwner(
		ctx,
		pastebins.GetManyInput{
			Page:        input.Page - 1,
			PerPage:     input.PerPage,
			OwnerUserID: input.OwnerUserID,
		},
	)
	if err != nil {
		return GetManyOutput{}, err
	}

	output := GetManyOutput{
		Items: make([]entity.Pastebin, 0, len(bins.Items)),
		Total: bins.Total,
	}

	for _, bin := range bins.Items {
		output.Items = append(output.Items, c.mapToEntity(bin))
	}

	for i, bin := range output.Items {
		if deleted, err := c.deleteIfNeed(ctx, bin); err != nil {
			return GetManyOutput{}, err
		} else if deleted {
			output.Items = append(output.Items[:i], output.Items[i+1:]...)
		}
	}

	return output, nil
}

func (c *Service) Delete(ctx context.Context, id string) error {
	err := c.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	if err := c.kv.Delete(ctx, makeKvStoreKey(id)); err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	return nil
}
