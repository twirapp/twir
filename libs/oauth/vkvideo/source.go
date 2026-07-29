package vkvideo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	goredis "github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/oauth"
	tokens "github.com/twirapp/twir/libs/repositories/tokens"
	vkvideobots "github.com/twirapp/twir/libs/repositories/vk_video_bots"
)

const singletonBotCredentialID = oauth.CredentialID("singleton")

type UserTokenSource interface {
	Token(context.Context, uuid.UUID) (oauth.Credential, error)
}

type SingletonBotTokenSource interface {
	Token(context.Context) (oauth.Credential, error)
}

type TransactionRunner interface {
	Do(context.Context, func(context.Context) error) error
}

type SourceOptions struct {
	ClientID          string
	ClientSecret      string
	RedirectURL       string
	APIBaseURL        string
	AuthBaseURL       string
	DevAPIBaseURL     string
	Redis             *goredis.Client
	CipherKey         string
	ClientFactory     ClientFactory
	TransactionRunner TransactionRunner
}

func (options SourceOptions) clientFactory() ClientFactory {
	if options.ClientFactory != nil {
		return options.ClientFactory
	}
	return func() (Client, error) {
		return vk.NewOAuthClient(vk.OAuthClientOpts{
			ClientID: options.ClientID, ClientSecret: options.ClientSecret, RedirectURL: options.RedirectURL,
			APIBaseURL: options.APIBaseURL, AuthBaseURL: options.AuthBaseURL, DevAPIBaseURL: options.DevAPIBaseURL,
		})
	}
}

func (options SourceOptions) locker() (*oauth.RedisLocker, error) {
	return oauth.NewRedisLocker(options.Redis, oauth.RedisLockerOptions{
		Prefix: "oauth", TTL: 30 * time.Second, RenewEvery: 10 * time.Second, Timeout: 5 * time.Second,
	})
}

func NewUserTokenSource(options SourceOptions, repository tokens.Repository) (UserTokenSource, error) {
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create VK Video user token locker: %w", err)
	}
	runtime, err := oauth.NewRefreshRuntime(
		userStore{repository: repository, cipherKey: options.CipherKey},
		refresher{factory: options.clientFactory()},
		locker,
		oauth.RuntimeOptions{Skew: time.Minute},
	)
	if err != nil {
		return nil, fmt.Errorf("create VK Video user token runtime: %w", err)
	}
	return userSource{runtime: runtime}, nil
}

func NewSingletonBotTokenSource(options SourceOptions, repository vkvideobots.Repository) (SingletonBotTokenSource, error) {
	if options.TransactionRunner == nil {
		return nil, fmt.Errorf("create VK Video bot token source: transaction runner is required")
	}
	locker, err := options.locker()
	if err != nil {
		return nil, fmt.Errorf("create VK Video bot token locker: %w", err)
	}
	runtime, err := oauth.NewRefreshRuntime(
		singletonBotStore{repository: repository, cipherKey: options.CipherKey},
		refresher{factory: options.clientFactory()},
		locker,
		oauth.RuntimeOptions{Skew: time.Minute},
	)
	if err != nil {
		return nil, fmt.Errorf("create VK Video bot token runtime: %w", err)
	}
	return singletonBotSource{runtime: runtime, transactionRunner: options.TransactionRunner}, nil
}

type userSource struct{ runtime *oauth.RefreshRuntime }

func (source userSource) Token(ctx context.Context, userID uuid.UUID) (oauth.Credential, error) {
	return source.runtime.Refresh(ctx, oauth.CredentialKey{Provider: Provider, ID: oauth.CredentialID(userID.String())})
}

type singletonBotSource struct {
	runtime           *oauth.RefreshRuntime
	transactionRunner TransactionRunner
}

func (source singletonBotSource) Token(ctx context.Context) (credential oauth.Credential, err error) {
	err = source.transactionRunner.Do(ctx, func(txCtx context.Context) error {
		credential, err = source.runtime.Refresh(txCtx, oauth.CredentialKey{Provider: Provider, ID: singletonBotCredentialID})
		return err
	})
	if err != nil {
		return oauth.Credential{}, fmt.Errorf("refresh VK Video bot token: %w", err)
	}
	return credential, nil
}
