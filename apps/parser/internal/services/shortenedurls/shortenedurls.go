package shortenedurls

import (
	"context"
	"errors"
	"net/url"
	"regexp"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/twirapp/twir/apps/parser/locales"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/i18n"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
)

type Opts struct {
	Repository shortenedurlsrepository.Repository
	Config     config.Config
}

func New(opts Opts) *Service {
	return &Service{
		repo:   opts.Repository,
		config: opts.Config,
	}
}

type Service struct {
	repo   shortenedurlsrepository.Repository
	config config.Config
}

type Link struct {
	Long  string
	Short string
}

var urlRegexp = regexp.MustCompile(`^https?://.*`)

func (c *Service) FindOrCreate(ctx context.Context, uri, actorId string) (*Link, error) {
	if !urlRegexp.MatchString(uri) {
		return nil, errors.New(i18n.GetCtx(ctx, locales.Translations.Services.Shortenedurls.Errors.InvalidUrl))
	}

	link, err := c.repo.GetByUrl(ctx, uri)
	if err != nil && !errors.Is(err, shortenedurlsrepository.ErrNotFound) {
		return nil, err
	}

	siteBaseUrl, _ := url.Parse(c.config.SiteBaseUrl)

	if link != model.Nil {
		siteBaseUrl.Path = "/s/" + link.ShortID
		return &Link{
			Long:  link.URL,
			Short: siteBaseUrl.String(),
		}, nil
	}

	newLink, err := c.repo.Create(
		ctx,
		shortenedurlsrepository.CreateInput{
			ShortID:         genId(),
			URL:             uri,
			CreatedByUserID: &actorId,
		},
	)
	if err != nil {
		return nil, err
	}

	siteBaseUrl.Path = "/s/" + newLink.ShortID
	return &Link{
		Long:  newLink.URL,
		Short: siteBaseUrl.String(),
	}, nil
}

var idAlphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func genId() string {
	id, _ := gonanoid.Generate(idAlphabet, 5)
	return id
}
