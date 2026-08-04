package mcp_oauth

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/google/uuid"
	appentity "github.com/twirapp/twir/apps/api-gql/internal/entity"
	dashboardaccess "github.com/twirapp/twir/apps/api-gql/internal/services/dashboard_access"
	"github.com/twirapp/twir/apps/api-gql/internal/services/users"
	cfg "github.com/twirapp/twir/libs/config"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	repository "github.com/twirapp/twir/libs/repositories/mcp_oauth"
	channelservice "github.com/twirapp/twir/libs/services/channels"
)

const manageBotSettings = "MANAGE_BOT_SETTINGS"

type Clock interface{ Now() time.Time }
type systemClock struct{}

func (systemClock) Now() time.Time { return time.Now() }

type DashboardSubject = dashboardaccess.Subject
type dashboardAccess interface {
	CanAccess(context.Context, DashboardSubject, uuid.UUID, string) (bool, error)
}
type userReader interface {
	GetByID(context.Context, string) (appentity.User, error)
}
type channelReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error)
}

type Dependencies struct {
	Repository      repository.Repository
	Users           userReader
	Channels        channelReader
	DashboardAccess dashboardAccess
	SiteBaseURL     string
	Clock           Clock
	Random          io.Reader
}
type Service struct {
	repository repository.Repository
	users      userReader
	channels   channelReader
	access     dashboardAccess
	resource   string
	metadata   string
	clock      Clock
	random     io.Reader
}

func NewFx(
	repo repository.Repository,
	usersService *users.Service,
	channelsService *channelservice.ChannelService,
	dashboardAccessService *dashboardaccess.Service,
	config cfg.Config,
) (*Service, error) {
	return New(Dependencies{
		Repository:      repo,
		Users:           usersService,
		Channels:        channelsService,
		DashboardAccess: dashboardAccessService,
		SiteBaseURL:     config.SiteBaseUrl,
		Clock:           systemClock{},
		Random:          cryptorand.Reader,
	})
}
func New(deps Dependencies) (*Service, error) {
	if deps.Repository == nil || deps.Users == nil || deps.Channels == nil || deps.DashboardAccess == nil || deps.Clock == nil || deps.Random == nil {
		return nil, fmt.Errorf("MCP OAuth service dependencies are required")
	}
	u, err := url.Parse(deps.SiteBaseURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("invalid MCP OAuth site base URL")
	}
	origin := &url.URL{Scheme: u.Scheme, Host: u.Host}
	return &Service{repository: deps.Repository, users: deps.Users, channels: deps.Channels, access: deps.DashboardAccess, resource: origin.JoinPath("api", "mcp").String(), metadata: origin.JoinPath(".well-known", "oauth-protected-resource", "api", "mcp").String(), clock: deps.Clock, random: deps.Random}, nil
}

func (s *Service) ProtectedResourceMetadataURL() string {
	return s.metadata
}

func (s *Service) resourceOrDefault(resource string) string {
	if resource == "" {
		return s.resource
	}
	return resource
}
