package shortenedurls

import (
	"context"
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository shortenedurlsrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository: opts.Repository,
	}
}

var idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genId() string {
	id, _ := gonanoid.Generate(idAlphabet, 5)
	return id
}

type Service struct {
	repository shortenedurlsrepository.Repository
}

type CreateInput struct {
	CreatedByUserID *string
	ShortID         string
	URL             string
}

func (c *Service) Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error) {
	shortId := input.ShortID
	if input.ShortID == "" {
		shortId = genId()
	}

	return c.repository.Create(
		ctx,
		shortenedurlsrepository.CreateInput{
			ShortID:         shortId,
			URL:             input.URL,
			CreatedByUserID: input.CreatedByUserID,
		},
	)
}

func (c *Service) GetByShortID(ctx context.Context, id string) (model.ShortenedUrl, error) {
	link, err := c.repository.GetByShortID(ctx, id)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrNotFound) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return link, nil
}

func (c *Service) GetByUrl(ctx context.Context, url string) (model.ShortenedUrl, error) {
	link, err := c.repository.GetByUrl(ctx, url)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrNotFound) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return link, nil
}
