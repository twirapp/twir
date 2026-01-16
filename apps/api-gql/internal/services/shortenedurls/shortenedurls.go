package shortenedurls

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository      shortenedurlsrepository.Repository
	ViewsRepository shortlinksviewsrepository.Repository
	WsRouter        wsrouter.WsRouter
}

func New(opts Opts) *Service {
	return &Service{
		repository:      opts.Repository,
		viewsRepository: opts.ViewsRepository,
		wsRouter:        opts.WsRouter,
		viewsSubs:       make(map[string]struct{}),
	}
}

var idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genId() string {
	id, _ := gonanoid.Generate(idAlphabet, 5)
	return id
}

const (
	shortLinkViewsSubscriptionKey = "api.shortLinkViews"
)

func shortLinkViewsSubscriptionKeyCreate(shortLinkId string) string {
	return shortLinkViewsSubscriptionKey + "." + shortLinkId
}

type Service struct {
	repository      shortenedurlsrepository.Repository
	viewsRepository shortlinksviewsrepository.Repository
	wsRouter        wsrouter.WsRouter
	viewsSubs       map[string]struct{}
	viewsSubsMu     sync.RWMutex
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
	Country     *string
	City        *string
}

func (c *Service) RecordView(ctx context.Context, input RecordViewInput) error {
	return c.viewsRepository.Create(
		ctx,
		shortlinksviewsrepository.CreateInput{
			ShortLinkID: input.ShortLinkID,
			UserID:      input.UserID,
			IP:          input.IP,
			UserAgent:   input.UserAgent,
			Country:     input.Country,
			City:        input.City,
			CreatedAt:   time.Now(),
		},
	)
}

type GetStatisticsInput struct {
	ShortLinkID string
	From        time.Time
	To          time.Time
	Interval    string
}

type StatisticsPoint struct {
	Timestamp int64 `json:"timestamp"`
	Count     int64 `json:"count"`
}

func (c *Service) GetStatistics(
	ctx context.Context,
	input GetStatisticsInput,
) ([]StatisticsPoint, error) {
	points, err := c.viewsRepository.GetStatistics(
		ctx,
		shortlinksviewsrepository.GetStatisticsInput{
			ShortLinkID: input.ShortLinkID,
			From:        input.From,
			To:          input.To,
			Interval:    input.Interval,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]StatisticsPoint, len(points))
	for i, p := range points {
		result[i] = StatisticsPoint{
			Timestamp: int64(p.Timestamp),
			Count:     int64(p.Count),
		}
	}

	return result, nil
}

type ShortLinkViewUpdate struct {
	ShortLinkID string
	TotalViews  int
	LastView    *shortlinksviewsrepository.View
}

func (c *Service) PublishViewUpdate(shortLinkID string, totalViews int, lastView *shortlinksviewsrepository.View) error {
	c.viewsSubsMu.RLock()
	defer c.viewsSubsMu.RUnlock()

	if _, ok := c.viewsSubs[shortLinkID]; !ok {
		return nil
	}

	update := ShortLinkViewUpdate{
		ShortLinkID: shortLinkID,
		TotalViews:  totalViews,
		LastView:    lastView,
	}

	return c.wsRouter.Publish(shortLinkViewsSubscriptionKeyCreate(shortLinkID), update)
}

func (c *Service) SubscribeToViewUpdates(ctx context.Context, shortLinkID string) <-chan ShortLinkViewUpdate {
	c.viewsSubsMu.Lock()
	c.viewsSubs[shortLinkID] = struct{}{}
	c.viewsSubsMu.Unlock()

	channel := make(chan ShortLinkViewUpdate)
	go func() {
		sub, err := c.wsRouter.Subscribe(
			[]string{
				shortLinkViewsSubscriptionKeyCreate(shortLinkID),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			c.viewsSubsMu.Lock()
			delete(c.viewsSubs, shortLinkID)
			c.viewsSubsMu.Unlock()
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-sub.GetChannel():
				if !ok {
					return
				}

				var update ShortLinkViewUpdate
				if err := json.Unmarshal(data, &update); err != nil {
					panic(err)
				}

				channel <- update
			}
		}
	}()

	return channel
}

type GetViewsInput struct {
	ShortLinkID string
	Page        int
	PerPage     int
}

type ViewOutput struct {
	ShortLinkID string
	UserID      *string
	Country     *string
	City        *string
	CreatedAt   time.Time
}

type GetViewsOutput struct {
	Views []ViewOutput
	Total int
}

func (c *Service) GetViews(ctx context.Context, input GetViewsInput) (GetViewsOutput, error) {
	result, err := c.viewsRepository.GetViews(
		ctx,
		shortlinksviewsrepository.GetViewsInput{
			ShortLinkID: input.ShortLinkID,
			Page:        input.Page,
			PerPage:     input.PerPage,
		},
	)
	if err != nil {
		return GetViewsOutput{}, err
	}

	views := make([]ViewOutput, len(result.Views))
	for i, v := range result.Views {
		views[i] = ViewOutput{
			ShortLinkID: v.ShortLinkID,
			UserID:      v.UserID,
			Country:     v.Country,
			City:        v.City,
			CreatedAt:   v.CreatedAt,
		}
	}

	return GetViewsOutput{
		Views: views,
		Total: result.Total,
	}, nil
}

type GetTopCountriesInput struct {
	ShortLinkID string
	Limit       int
}

type CountryStatsOutput struct {
	Country string
	Count   uint64
}

func (c *Service) GetTopCountries(ctx context.Context, input GetTopCountriesInput) ([]CountryStatsOutput, error) {
	result, err := c.viewsRepository.GetTopCountries(
		ctx,
		shortlinksviewsrepository.GetTopCountriesInput{
			ShortLinkID: input.ShortLinkID,
			Limit:       input.Limit,
		},
	)
	if err != nil {
		return nil, err
	}

	output := make([]CountryStatsOutput, len(result))
	for i, cs := range result {
		output[i] = CountryStatsOutput{
			Country: cs.Country,
			Count:   cs.Count,
		}
	}

	return output, nil
}
