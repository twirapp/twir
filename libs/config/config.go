package cfg

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type ValorantConfig struct {
	HenrikApiKey string `required:"false" envconfig:"VALORANT_HENRIK_API_KEY"`
	ClientID     string `required:"false" envconfig:"VALORANT_CLIENT_ID"`
	ClientSecret string `required:"false" envconfig:"VALORANT_CLIENT_SECRET"`
	RedirectURL  string `required:"false" envconfig:"VALORANT_REDIRECT_URL"`
	RiotApiKey   string `required:"false" envconfig:"VALORANT_RIOT_API_KEY"`
}

type LastFMConfig struct {
	ApiKey       string `required:"false" envconfig:"LASTFM_API_KEY"`
	ClientSecret string `required:"false" envconfig:"LASTFM_CLIENT_SECRET"`
	RedirectURL  string `required:"false" envconfig:"LASTFM_REDIRECT_URL"`
}

type Config struct {
	BotAccessToken  string `required:"false"    envconfig:"BOT_ACCESS_TOKEN"`
	BotRefreshToken string `required:"false"    envconfig:"BOT_REFRESH_TOKEN"`

	TrustedProxies []string `envconfig:"TRUSTED_PROXIES"`

	RedisUrl             string `required:"true"  default:"redis://localhost:6379/0"    envconfig:"REDIS_URL"`
	TwitchClientId       string `required:"true"                                        envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret   string `required:"true"                                        envconfig:"TWITCH_CLIENTSECRET"`
	KickClientId         string `required:"false"                                       envconfig:"KICK_CLIENT_ID"`
	KickClientSecret     string `required:"false"                                       envconfig:"KICK_CLIENT_SECRET"`
	DatabaseUrl          string `required:"true"                                        envconfig:"DATABASE_URL"`
	MigrationDatabaseUrl string `required:"false"                                       envconfig:"MIGRATION_DATABASE_URL"`
	ClickhouseUrl        string `required:"true"  default:"clickhouse://twir:twir@127.0.0.1:9000/twir" envconfig:"CLICKHOUSE_URL"`
	AppEnv               string `required:"true"  default:"development"                 envconfig:"APP_ENV"`
	SentryDsn            string `required:"false"                                       envconfig:"SENTRY_DSN"`
	SiteBaseUrl          string `required:"false" default:"http://localhost:3005" envconfig:"SITE_BASE_URL"`
	TokensCipherKey      string `required:"false" default:"pnyfwfiulmnqlhkvixaeligpprcnlyke" envconfig:"TOKENS_CIPHER_KEY"`
	TTSServiceUrl        string `required:"false" default:"localhost:7001" envconfig:"TTS_SERVICE_URL"`
	OdesliApiKey         string `required:"false" envconfig:"ODESLI_API_KEY"`

	S3PublicUrl   string `required:"false" envconfig:"CDN_PUBLIC_URL"`
	S3Host        string `required:"false" envconfig:"CDN_HOST"`
	S3Bucket      string `required:"false" envconfig:"CDN_BUCKET"`
	S3Region      string `required:"false" envconfig:"CDN_REGION"`
	S3AccessToken string `required:"false" envconfig:"CDN_ACCESS_TOKEN"`
	S3SecretToken string `required:"false" envconfig:"CDN_SECRET_TOKEN"`

	UploaderMaxFileSizeBytes     int64         `required:"false" default:"26214400" envconfig:"UPLOADER_MAX_FILE_SIZE_BYTES"`
	UploaderAnonymousFileTTL     time.Duration `required:"false" default:"720h" envconfig:"UPLOADER_ANONYMOUS_FILE_TTL"`
	UploaderAuthenticatedFileTTL time.Duration `required:"false" default:"4320h" envconfig:"UPLOADER_AUTHENTICATED_FILE_TTL"`
	UploaderS3Host               string        `required:"false" envconfig:"UPLOADER_S3_HOST"`
	UploaderS3Bucket             string        `required:"false" default:"uploader-s3" envconfig:"UPLOADER_S3_BUCKET"`
	UploaderS3Region             string        `required:"false" envconfig:"UPLOADER_S3_REGION"`
	UploaderS3AccessToken        string        `required:"false" envconfig:"UPLOADER_S3_ACCESS_TOKEN"`
	UploaderS3SecretToken        string        `required:"false" envconfig:"UPLOADER_S3_SECRET_TOKEN"`
	UploaderS3PublicURL          string        `required:"false" envconfig:"UPLOADER_S3_PUBLIC_URL"`

	DiscordClientID                        string `required:"false" envconfig:"DISCORD_CLIENT_ID"`
	DiscordClientSecret                    string `required:"false" envconfig:"DISCORD_CLIENT_SECRET"`
	DiscordBotToken                        string `required:"false" envconfig:"DISCORD_BOT_TOKEN"`
	DiscordFeedbackUrl                     string `required:"false" envconfig:"DISCORD_FEEDBACK_URL"`
	DiscordNotificationsChannelID          string `required:"false" envconfig:"DISCORD_NOTIFICATIONS_CHANNEL_ID"`
	DiscordNotificationsHistoryLimit       uint   `required:"false" default:"100" envconfig:"DISCORD_NOTIFICATIONS_HISTORY_LIMIT"`
	DiscordNotificationsMaxAttachmentBytes int64  `required:"false" default:"26214400" envconfig:"DISCORD_NOTIFICATIONS_MAX_ATTACHMENT_BYTES"`
	GithubWebhookSecret                    string `required:"false" envconfig:"GITHUB_WEBHOOK_SECRET"`

	OpenWeatherMapApiKey string `required:"false" envconfig:"OPENWEATHERMAP_API_KEY"`

	TemporalHost string `required:"false" default:"localhost:7233" envconfig:"TEMPORAL_HOST"`

	SevenTvToken string `required:"false" envconfig:"SEVENTV_TOKEN"`

	// OpenTelemetry configuration
	OtelEndpoint       string `required:"false" envconfig:"OTEL_ENDPOINT" default:"localhost:4317"`
	OtelHeaders        string `required:"false" envconfig:"OTEL_HEADERS"`
	OtelInsecure       bool   `required:"false" default:"true" envconfig:"OTEL_INSECURE"`
	OtelTracingEnabled bool   `required:"false" default:"true" envconfig:"OTEL_TRACING_ENABLED"`
	OtelMetricsEnabled bool   `required:"false" default:"true" envconfig:"OTEL_METRICS_ENABLED"`

	NatsUrl string `required:"false" default:"localhost:4222" envconfig:"NATS_URL"`

	Valorant ValorantConfig
	LastFM   LastFMConfig

	ToxicityAddr        string `required:"false" envconfig:"TOXICITY_ADDR"`
	MusicRecognizerAddr string `required:"false" envconfig:"MUSIC_RECOGNIZER_ADDR"`

	StreamElementsClientId     string `required:"false" envconfig:"STREAM_ELEMENTS_CLIENT_ID"`
	StreamElementsClientSecret string `required:"false" envconfig:"STREAM_ELEMENTS_CLIENT_SECRET"`

	SpotifyClientID string `required:"false" envconfig:"SPOTIFY_CLIENT_ID"`
	SpotifySecret   string `required:"false" envconfig:"SPOTIFY_CLIENT_SECRET"`

	ExecutronAddr           string `required:"false" default:"http://localhost:7003" envconfig:"EXECUTRON_ADDR"`
	ExecutronCfClientId     string `required:"false" envconfig:"EXECUTRON_CF_CLIENT_ID"`
	ExecutronCfClientSecret string `required:"false" envconfig:"EXECUTRON_CF_CLIENT_SECRET"`

	EventSubDisableSignatureVerification bool   `required:"false" default:"false" envconfig:"EVENTSUB_DISABLE_SIGNATURE_VERIFICATION"`
	EventsubHttpPort                     int    `required:"false" default:"3030"  envconfig:"EVENTSUB_HTTP_PORT"`
	EventSubCallbackBaseUrl              string `required:"false" envconfig:"EVENTSUB_CALLBACK_BASE_URL"`

	DonationAlertsClientId string `required:"false" envconfig:"DONATIONALERTS_CLIENT_ID"`
	DonationAlertsSecret   string `required:"false" envconfig:"DONATIONALERTS_CLIENT_SECRET"`

	VKClientId       string `required:"false" envconfig:"VK_CLIENT_ID"`
	VKClientSecret   string `required:"false" envconfig:"VK_CLIENT_SECRET"`
	VkAppAccessToken string `required:"false" envconfig:"VK_APP_ACCESS_TOKEN"`

	VKVideoClientID      string `required:"false" envconfig:"VK_VIDEO_CLIENT_ID"`
	VKVideoClientSecret  string `required:"false" envconfig:"VK_VIDEO_CLIENT_SECRET"`
	VKVideoWebhookSecret string `required:"false" envconfig:"VK_VIDEO_WEBHOOK_SECRET"`
	VKVideoAPIBaseURL    string `required:"false" default:"https://api.live.vkvideo.ru" envconfig:"VK_VIDEO_API_BASE_URL"`
	VKVideoAuthBaseURL   string `required:"false" default:"https://auth.live.vkvideo.ru" envconfig:"VK_VIDEO_AUTH_BASE_URL"`
	VKVideoDevAPIBaseURL string `required:"false" default:"https://apidev.live.vkvideo.ru" envconfig:"VK_VIDEO_DEVAPI_BASE_URL"`

	// VkProxyUrl is an optional HTTP/SOCKS5 proxy URL used for the VK Video
	// Live Centrifugo (pubsub) WebSocket connection, e.g. "http://host:port"
	// or "socks5://user:pass@host:port".
	VkProxyUrl string `required:"false" envconfig:"VK_PROXY_URL"`

	YouTubeClientID     string `required:"false" envconfig:"YOUTUBE_CLIENT_ID"`
	YouTubeClientSecret string `required:"false" envconfig:"YOUTUBE_CLIENT_SECRET"`

	FaceitClientId     string `required:"false" envconfig:"FACEIT_CLIENT_ID"`
	FaceitClientSecret string `required:"false" envconfig:"FACEIT_CLIENT_SECRET"`
	FaceitApiKey       string `required:"false" envconfig:"FACEIT_API_KEY"`

	DeeplApiKey                       string `required:"false" envconfig:"DEEPL_API_KEY"`
	GoogleTranslateServiceAccountJson string `required:"false" envconfig:"GOOGLE_TRANSLATE_SERVICE_ACCOUNT_JSON"`

	StreamlabsClientId     string `required:"false" envconfig:"STREAMLABS_CLIENT_ID"`
	TwitchMockEnabled      bool   `required:"false" default:"false" envconfig:"TWITCH_MOCK_ENABLED"`
	TwitchMockApiUrl       string `required:"false" default:"http://twitch-mock:7777/helix" envconfig:"TWITCH_MOCK_API_URL"`
	TwitchMockAuthUrl      string `required:"false" default:"http://twitch-mock:7777" envconfig:"TWITCH_MOCK_AUTH_URL"`
	TwitchMockWsUrl        string `required:"false" default:"ws://twitch-mock:8081/ws" envconfig:"TWITCH_MOCK_WS_URL"`
	StreamlabsClientSecret string `required:"false" envconfig:"STREAMLABS_CLIENT_SECRET"`

	SecretsEncryptionKey string `required:"false" default:"0123456789abcdef0123456789abcdef" envconfig:"SECRETS_ENCRYPTION_KEY"`
}

func (c Config) BuildS3PublicURL(objectPath string) string {
	parts := []string{strings.TrimRight(c.S3PublicUrl, "/")}
	if c.AppEnv == "development" {
		parts = append(parts, strings.Trim(c.S3Bucket, "/"))
	}
	if objectPath = strings.TrimLeft(objectPath, "/"); objectPath != "" {
		parts = append(parts, objectPath)
	}
	return strings.Join(parts, "/")
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func (c *Config) IsDevelopment() bool {
	return c.AppEnv == "development"
}

func (c *Config) IsVkVideoEnabled() bool {
	return c.VKVideoClientID != "" && c.VKVideoClientSecret != ""
}

func (c *Config) IsYouTubeEnabled() bool {
	return c.YouTubeClientID != "" && c.YouTubeClientSecret != ""
}

func (c *Config) GetTwitchCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("login").String()
}

func (c *Config) GetVkCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("login", "vk").String()
}

func (c *Config) GetVkVideoBotCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("api", "auth", "vk-video", "bot-callback").String()
}

func (c *Config) GetKickCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("login", "kick").String()
}

func (c *Config) GetYouTubeCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("login", "youtube").String()
}

func (c *Config) GetYouTubeBotCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("api", "auth", "youtube", "bot-callback").String()
}

func NewWithEnvPath(envPath string) (*Config, error) {
	var newCfg Config
	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}
	if err := newCfg.validateVKVideo(); err != nil {
		return nil, err
	}

	return &newCfg, nil
}

func (c *Config) validateVKVideo() error {
	if !c.IsVkVideoEnabled() {
		return nil
	}

	for name, value := range map[string]string{
		"VK_VIDEO_CLIENT_ID":      c.VKVideoClientID,
		"VK_VIDEO_CLIENT_SECRET":  c.VKVideoClientSecret,
		"VK_VIDEO_WEBHOOK_SECRET": c.VKVideoWebhookSecret,
	} {
		if strings.TrimSpace(value) == "" {
			return fmt.Errorf("%s is required when VK_VIDEO_ENABLED is true", name)
		}
	}
	return nil
}

func New() (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(wd, "/workspace") {
		wd = "/workspace"
	} else {
		wd = filepath.Join(wd, "..", "..")
	}

	envPath := filepath.Join(wd, ".env")

	return NewWithEnvPath(envPath)
}

func NewFx() Config {
	config, err := New()
	if err != nil {
		panic(err)
	}

	return *config
}

func NewFxWithPath(path string) Config {
	config, err := NewWithEnvPath(path)
	if err != nil {
		panic(err)
	}

	return *config
}
