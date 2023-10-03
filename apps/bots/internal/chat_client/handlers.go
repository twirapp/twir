package chat_client

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	model "github.com/satont/twir/libs/gomodels"
)

type BotInstance struct {
	BotClient *ChatClient
	Db        *model.Bots
}

var greetingsCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "bots_greetings_counter",
		Help: "The total number of processed greetings",
		// ConstLabels: labels,
	},
)
var keywordsCounter = promauto.NewCounter(
	prometheus.CounterOpts{
		Name: "bots_keywords_counter",
		Help: "The total number of processed keywords",
		// ConstLabels: labels,
	},
)
