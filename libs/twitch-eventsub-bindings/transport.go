package eventsub_bindings

import "encoding/json"

type Pagination struct {
	Cursor string `json:"cursor"`
}

type Request struct {
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Condition interface{} `json:"condition"`
	Transport Transport   `json:"transport"`
}

type RequestStatus struct {
	Data         []Subscription `json:"data"`
	Total        int            `json:"total"`
	TotalCost    int            `json:"total_cost"`
	MaxTotalCost int            `json:"max_total_cost"`
	Pagination   *Pagination    `json:"pagination"`
}

type EventNotification struct {
	Subscription Subscription    `json:"subscription"`
	Event        json.RawMessage `json:"event"`
}

type SubscriptionChallenge struct {
	Challenge    string       `json:"challenge"`
	Subscription Subscription `json:"subscription"`
}

type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

type ResponseHeaders struct {
	MessageID           string `header:"Twitch-Eventsub-Message-Id"`
	MessageRetry        int    `header:"Twitch-Eventsub-Message-Retry"`
	MessageType         string `header:"Twitch-Eventsub-Message-Type"`
	MessageSignature    string `header:"Twitch-Eventsub-Message-Signature"`
	MessageTimestamp    string `header:"Twitch-Eventsub-Message-Timestamp"`
	SubscriptionType    string `header:"Twitch-Eventsub-Subscription-Type"`
	SubscriptionVersion string `header:"Twitch-Eventsub-Subscription-Version"`
}
