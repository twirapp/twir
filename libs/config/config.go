package cfg

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BotAccessToken  string `required:"false"    envconfig:"BOT_ACCESS_TOKEN"`
	BotRefreshToken string `required:"false"    envconfig:"BOT_REFRESH_TOKEN"`

	RedisUrl           string `required:"true"  default:"redis://localhost:6379/0"    envconfig:"REDIS_URL"`
	TwitchClientId     string `required:"true"                                        envconfig:"TWITCH_CLIENTID"`
	TwitchClientSecret string `required:"true"                                        envconfig:"TWITCH_CLIENTSECRET"`
	DatabaseUrl        string `required:"true"                                        envconfig:"DATABASE_URL"`
	ClickhouseUrl      string `required:"true"  default:"clickhouse://twir:twir@127.0.0.1:9000/twir" envconfig:"CLICKHOUSE_URL"`
	AppEnv             string `required:"true"  default:"development"                 envconfig:"APP_ENV"`
	SentryDsn          string `required:"false"                                       envconfig:"SENTRY_DSN"`
	SiteBaseUrl        string `required:"false" default:"http://localhost:3005" envconfig:"SITE_BASE_URL"`
	TokensCipherKey    string `required:"false" default:"pnyfwfiulmnqlhkvixaeligpprcnlyke" envconfig:"TOKENS_CIPHER_KEY"`
	TTSServiceUrl      string `required:"false" default:"localhost:7001" envconfig:"TTS_SERVICE_URL"`
	OdesliApiKey       string `required:"false" envconfig:"ODESLI_API_KEY"`

	S3PublicUrl   string `required:"false" envconfig:"CDN_PUBLIC_URL"`
	S3Host        string `required:"false" envconfig:"CDN_HOST"`
	S3Bucket      string `required:"false" envconfig:"CDN_BUCKET"`
	S3Region      string `required:"false" envconfig:"CDN_REGION"`
	S3AccessToken string `required:"false" envconfig:"CDN_ACCESS_TOKEN"`
	S3SecretToken string `required:"false" envconfig:"CDN_SECRET_TOKEN"`

	DiscordClientID     string `required:"false" envconfig:"DISCORD_CLIENT_ID"`
	DiscordClientSecret string `required:"false" envconfig:"DISCORD_CLIENT_SECRET"`
	DiscordBotToken     string `required:"false" envconfig:"DISCORD_BOT_TOKEN"`
	DiscordFeedbackUrl  string `required:"false" envconfig:"DISCORD_FEEDBACK_URL"`

	OpenWeatherMapApiKey string `required:"false" envconfig:"OPENWEATHERMAP_API_KEY"`

	TemporalHost string `required:"false" default:"localhost:7233" envconfig:"TEMPORAL_HOST"`

	SevenTvToken string `required:"false" envconfig:"SEVENTV_TOKEN"`

	UptraceDsn           string `required:"false" envconfig:"UPTRACE_DSN"`
	NatsUrl              string `required:"false" default:"localhost:4222" envconfig:"NATS_URL"`
	ValorantHenrikApiKey string `required:"false" envconfig:"VALORANT_HENRIK_API_KEY"`

	ToxicityAddr        string `required:"false" envconfig:"TOXICITY_ADDR"`
	MusicRecognizerAddr string `required:"false" envconfig:"MUSIC_RECOGNIZER_ADDR"`

	StreamElementsClientId     string `required:"false" envconfig:"STREAM_ELEMENTS_CLIENT_ID"`
	StreamElementsClientSecret string `required:"false" envconfig:"STREAM_ELEMENTS_CLIENT_SECRET"`

	ExecutronAddr string `required:"false" default:"http://localhost:7003" envconfig:"EXECUTRON_ADDR"`

	EventSubDisableSignatureVerification bool `required:"false" default:"false" envconfig:"EVENTSUB_DISABLE_SIGNATURE_VERIFICATION"`
}

func (c *Config) GetTwitchCallbackUrl() string {
	u, err := url.Parse(c.SiteBaseUrl)
	if err != nil {
		panic(err)
	}

	return u.JoinPath("login").String()
}

func NewWithEnvPath(envPath string) (*Config, error) {
	var newCfg Config
	_ = godotenv.Load(envPath)

	if err := envconfig.Process("", &newCfg); err != nil {
		return nil, err
	}

	return &newCfg, nil
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
