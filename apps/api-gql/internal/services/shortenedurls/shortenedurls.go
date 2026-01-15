package shortenedurls

import (
	"context"
	"errors"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository      shortenedurlsrepository.Repository
	ViewsRepository shortlinksviewsrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository:      opts.Repository,
		viewsRepository: opts.ViewsRepository,
	}
}

var idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genId() string {
	id, _ := gonanoid.Generate(idAlphabet, 5)
	return id
}

type Service struct {
	repository      shortenedurlsrepository.Repository
	viewsRepository shortlinksviewsrepository.Repository
}

type CreateInput struct {
	CreatedByUserID *string
	ShortID         string
	URL             string
	UserIp          *string
	UserAgent       *string
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
			UserIp:          input.UserIp,
			UserAgent:       input.UserAgent,
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

type UpdateInput struct {
	Views *int
}

func (c *Service) Update(ctx context.Context, id string, input UpdateInput) error {
	_, err := c.repository.Update(
		ctx,
		id,
		shortenedurlsrepository.UpdateInput{
			Views: input.Views,
		},
	)
	return err
}

type GetListInput struct {
	Page        int
	PerPage     int
	OwnerUserID *string
}

type GetListOutput struct {
	List  []entity.ShortenedUrl
	Total int
}

func (c *Service) GetList(ctx context.Context, input GetListInput) (GetListOutput, error) {
	page := input.Page
	if page < 0 {
		page = 0
	}

	perPage := input.PerPage
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	data, err := c.repository.GetList(
		ctx, shortenedurlsrepository.GetListInput{
			Page:    page,
			PerPage: perPage,
			UserID:  input.OwnerUserID,
		},
	)
	if err != nil {
		return GetListOutput{}, err
	}

	converted := make([]entity.ShortenedUrl, 0, len(data.Items))
	for _, link := range data.Items {
		converted = append(
			converted,
			entity.ShortenedUrl{
				ID:          link.ShortID,
				Link:        link.URL,
				Views:       link.Views,
				CreatedAt:   link.CreatedAt,
				UpdatedAt:   link.UpdatedAt,
				OwnerUserID: link.CreatedByUserId,
			},
		)
	}

	return GetListOutput{
		List:  converted,
		Total: data.Total,
	}, nil
}

func (c *Service) Delete(ctx context.Context, id string) error {
	return c.repository.Delete(ctx, id)
}

func (c *Service) GetManyByShortIDs(ctx context.Context, ids []string) (
	[]model.ShortenedUrl,
	error,
) {
	return c.repository.GetManyByShortIDs(ctx, ids)
}

type RecordViewInput struct {
	ShortLinkID string
	UserID      *string
	IP          *string
	UserAgent   *string
}

func (c *Service) RecordView(ctx context.Context, input RecordViewInput) error {
	return c.viewsRepository.Create(
		ctx,
		shortlinksviewsrepository.CreateInput{
			ShortLinkID: input.ShortLinkID,
			UserID:      input.UserID,
			IP:          input.IP,
			UserAgent:   input.UserAgent,
			CreatedAt:   time.Now(),
		},
	)
}
