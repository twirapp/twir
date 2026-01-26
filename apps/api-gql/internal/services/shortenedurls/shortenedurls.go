package shortenedurls

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	shortlinkscustomdomainsrepo "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository              shortenedurlsrepository.Repository
	ViewsRepository         shortlinksviewsrepository.Repository
	CustomDomainsRepository shortlinkscustomdomainsrepo.Repository
	WsRouter                wsrouter.WsRouter
}

func New(opts Opts) *Service {
	return &Service{
		repository:              opts.Repository,
		viewsRepository:         opts.ViewsRepository,
		customDomainsRepository: opts.CustomDomainsRepository,
		wsRouter:                opts.WsRouter,
		viewsSubs:               make(map[string]struct{}),
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

func shortLinkViewsSubscriptionKeyCreate(viewKey string) string {
	return shortLinkViewsSubscriptionKey + "." + viewKey
}

const shortLinkDomainSeparator = "|"

func EncodeShortLinkKey(domain *string, shortID string) string {
	if domain == nil || *domain == "" {
		return shortID
	}

	return *domain + shortLinkDomainSeparator + shortID
}

func DecodeShortLinkKey(value string) (*string, string) {
	if value == "" {
		return nil, ""
	}

	parts := strings.SplitN(value, shortLinkDomainSeparator, 2)
	if len(parts) == 2 && parts[0] != "" {
		domain := parts[0]
		return &domain, parts[1]
	}

	return nil, value
}

type Service struct {
	repository              shortenedurlsrepository.Repository
	viewsRepository         shortlinksviewsrepository.Repository
	customDomainsRepository shortlinkscustomdomainsrepo.Repository
	wsRouter                wsrouter.WsRouter
	viewsSubs               map[string]struct{}
	viewsSubsMu             sync.RWMutex
}

func (c *Service) resolveDomainID(ctx context.Context, domain *string) (*string, error) {
	if domain == nil || *domain == "" {
		return nil, nil
	}

	customDomain, err := c.customDomainsRepository.GetByDomain(ctx, *domain)
	if err != nil {
		if errors.Is(err, shortlinkscustomdomainsrepo.ErrNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &customDomain.ID, nil
}

func (c *Service) resolveDomainIDRequired(ctx context.Context, domain *string) (*string, error) {
	domainID, err := c.resolveDomainID(ctx, domain)
	if err != nil {
		return nil, err
	}
	if domain != nil && domainID == nil {
		return nil, shortlinkscustomdomainsrepo.ErrNotFound
	}
	return domainID, nil
}

type CreateInput struct {
	CreatedByUserID *string
	ShortID         string
	URL             string
	UserIp          *string
	UserAgent       *string
	Domain          *string
}

func (c *Service) Create(ctx context.Context, input CreateInput) (model.ShortenedUrl, error) {
	domainID, err := c.resolveDomainIDRequired(ctx, input.Domain)
	if err != nil {
		return model.Nil, err
	}

	shortId := input.ShortID
	if input.ShortID == "" {
		shortId = genId()
	} else {
		existing, err := c.repository.GetByShortID(ctx, domainID, input.ShortID)
		if err != nil && !errors.Is(err, shortenedurlsrepository.ErrNotFound) {
			return model.Nil, err
		}
		if err == nil && !existing.IsNil() {
			return model.Nil, shortenedurlsrepository.ErrShortIDAlreadyExists
		}
	}

	return c.repository.Create(
		ctx,
		shortenedurlsrepository.CreateInput{
			ShortID:         shortId,
			URL:             input.URL,
			CreatedByUserID: input.CreatedByUserID,
			UserIp:          input.UserIp,
			UserAgent:       input.UserAgent,
			DomainID:        domainID,
		},
	)
}

func (c *Service) MoveLinksToDefaultDomain(
	ctx context.Context,
	userID string,
	domain string,
) error {
	if userID == "" || domain == "" {
		return errors.New("user ID and domain are required")
	}

	domainID, err := c.resolveDomainIDRequired(ctx, &domain)
	if err != nil {
		return err
	}

	conflicts, err := c.repository.CountDomainShortIDConflicts(ctx, *domainID, userID)
	if err != nil {
		return err
	}
	if conflicts > 0 {
		return shortenedurlsrepository.ErrShortIDAlreadyExists
	}

	return c.repository.ClearDomainForUser(ctx, *domainID, userID)
}

func (c *Service) GetByShortID(
	ctx context.Context,
	domain *string,
	id string,
) (model.ShortenedUrl, error) {
	domainID, err := c.resolveDomainID(ctx, domain)
	if err != nil {
		return model.Nil, err
	}
	if domain != nil && domainID == nil {
		return model.Nil, nil
	}

	link, err := c.repository.GetByShortID(ctx, domainID, id)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrNotFound) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return link, nil
}

func (c *Service) GetByUrl(
	ctx context.Context,
	domain *string,
	url string,
) (model.ShortenedUrl, error) {
	domainID, err := c.resolveDomainID(ctx, domain)
	if err != nil {
		return model.Nil, err
	}
	if domain != nil && domainID == nil {
		return model.Nil, nil
	}

	link, err := c.repository.GetByUrl(ctx, domainID, url)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrNotFound) {
			return model.Nil, nil
		}
		return model.Nil, err
	}

	return link, nil
}

type UpdateInput struct {
	Views       *int
	ShortID     *string
	URL         *string
	Domain      *string
	ClearDomain bool
}

func (c *Service) Update(
	ctx context.Context,
	domain *string,
	id string,
	input UpdateInput,
) (entity.ShortenedUrl, error) {
	domainID, err := c.resolveDomainID(ctx, domain)
	if err != nil {
		return entity.Nil, err
	}
	if domain != nil && domainID == nil {
		return entity.Nil, shortenedurlsrepository.ErrNotFound
	}

	var targetDomainID *string
	if input.Domain != nil {
		targetDomainID, err = c.resolveDomainIDRequired(ctx, input.Domain)
		if err != nil {
			return entity.Nil, err
		}
	}

	updatedModel, err := c.repository.Update(
		ctx,
		domainID,
		id,
		shortenedurlsrepository.UpdateInput{
			Views:       input.Views,
			ShortID:     input.ShortID,
			URL:         input.URL,
			DomainID:    targetDomainID,
			ClearDomain: input.ClearDomain,
		},
	)
	if err != nil {
		return entity.Nil, err
	}

	return entity.ShortenedUrl{
		ID:          updatedModel.ShortID,
		Link:        updatedModel.URL,
		Views:       updatedModel.Views,
		CreatedAt:   updatedModel.CreatedAt,
		UpdatedAt:   updatedModel.UpdatedAt,
		OwnerUserID: updatedModel.CreatedByUserId,
		Domain:      updatedModel.Domain,
	}, nil
}

type GetListInput struct {
	Page        int
	PerPage     int
	OwnerUserID *string
	SortBy      string // "views" or "created_at"
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
			SortBy:  input.SortBy,
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
				Domain:      link.Domain,
			},
		)
	}

	return GetListOutput{
		List:  converted,
		Total: data.Total,
	}, nil
}

func (c *Service) Delete(ctx context.Context, domain *string, id string) error {
	domainID, err := c.resolveDomainID(ctx, domain)
	if err != nil {
		return err
	}
	if domain != nil && domainID == nil {
		return nil
	}

	return c.repository.Delete(ctx, domainID, id)
}

func (c *Service) GetManyByShortIDs(ctx context.Context, ids []string) ([]model.ShortenedUrl, error) {
	if len(ids) == 0 {
		return []model.ShortenedUrl{}, nil
	}

	grouped := make(map[string][]string)
	var defaultDomainIDs []string

	for _, rawID := range ids {
		domain, shortID := DecodeShortLinkKey(rawID)
		if shortID == "" {
			continue
		}
		if domain == nil {
			defaultDomainIDs = append(defaultDomainIDs, shortID)
			continue
		}
		grouped[*domain] = append(grouped[*domain], shortID)
	}

	results := make([]model.ShortenedUrl, 0, len(ids))
	if len(defaultDomainIDs) > 0 {
		links, err := c.repository.GetManyByShortIDs(ctx, nil, defaultDomainIDs)
		if err != nil {
			return nil, err
		}
		results = append(results, links...)
	}

	for domain, domainIDs := range grouped {
		if len(domainIDs) == 0 {
			continue
		}
		domainValue := domain
		domainID, err := c.resolveDomainID(ctx, &domainValue)
		if err != nil {
			return nil, err
		}
		if domainID == nil {
			continue
		}
		links, err := c.repository.GetManyByShortIDs(ctx, domainID, domainIDs)
		if err != nil {
			return nil, err
		}
		results = append(results, links...)
	}

	return results, nil
}

type RecordViewInput struct {
	ShortLinkID string
	Domain      *string
	UserID      *string
	IP          *string
	UserAgent   *string
	Country     *string
	City        *string
}

func (c *Service) RecordView(ctx context.Context, input RecordViewInput) error {
	viewKey := EncodeShortLinkKey(input.Domain, input.ShortLinkID)
	return c.viewsRepository.Create(
		ctx,
		shortlinksviewsrepository.CreateInput{
			ShortLinkID: viewKey,
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
	Domain      *string
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
	viewKey := EncodeShortLinkKey(input.Domain, input.ShortLinkID)
	points, err := c.viewsRepository.GetStatistics(
		ctx,
		shortlinksviewsrepository.GetStatisticsInput{
			ShortLinkID: viewKey,
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

func (c *Service) PublishViewUpdate(
	domain *string,
	shortLinkID string,
	totalViews int,
	lastView *shortlinksviewsrepository.View,
) error {
	viewKey := EncodeShortLinkKey(domain, shortLinkID)
	c.viewsSubsMu.RLock()
	defer c.viewsSubsMu.RUnlock()

	if _, ok := c.viewsSubs[viewKey]; !ok {
		return nil
	}

	update := ShortLinkViewUpdate{
		ShortLinkID: shortLinkID,
		TotalViews:  totalViews,
		LastView:    lastView,
	}

	return c.wsRouter.Publish(shortLinkViewsSubscriptionKeyCreate(viewKey), update)
}

func (c *Service) SubscribeToViewUpdates(
	ctx context.Context,
	domain *string,
	shortLinkID string,
) <-chan ShortLinkViewUpdate {
	viewKey := EncodeShortLinkKey(domain, shortLinkID)
	c.viewsSubsMu.Lock()
	c.viewsSubs[viewKey] = struct{}{}
	c.viewsSubsMu.Unlock()

	channel := make(chan ShortLinkViewUpdate)
	go func() {
		sub, err := c.wsRouter.Subscribe(
			[]string{
				shortLinkViewsSubscriptionKeyCreate(viewKey),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			c.viewsSubsMu.Lock()
			delete(c.viewsSubs, viewKey)
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
	Domain      *string
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
	viewKey := EncodeShortLinkKey(input.Domain, input.ShortLinkID)
	result, err := c.viewsRepository.GetViews(
		ctx,
		shortlinksviewsrepository.GetViewsInput{
			ShortLinkID: viewKey,
			Page:        input.Page,
			PerPage:     input.PerPage,
		},
	)
	if err != nil {
		return GetViewsOutput{}, err
	}

	views := make([]ViewOutput, len(result.Views))
	for i, v := range result.Views {
		_, shortID := DecodeShortLinkKey(v.ShortLinkID)
		views[i] = ViewOutput{
			ShortLinkID: shortID,
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
	Domain      *string
	Limit       int
}

type CountryStatsOutput struct {
	Country string
	Count   uint64
}

func (c *Service) GetTopCountries(ctx context.Context, input GetTopCountriesInput) ([]CountryStatsOutput, error) {
	viewKey := EncodeShortLinkKey(input.Domain, input.ShortLinkID)
	result, err := c.viewsRepository.GetTopCountries(
		ctx,
		shortlinksviewsrepository.GetTopCountriesInput{
			ShortLinkID: viewKey,
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
