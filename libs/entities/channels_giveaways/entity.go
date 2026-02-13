package channels_giveaways

import (
	"time"

	"github.com/google/uuid"
)

type GiveawayType string

const (
	GiveawayTypeKeyword        GiveawayType = "KEYWORD"
	GiveawayTypeOnlineChatters GiveawayType = "ONLINE_CHATTERS"
)

type Giveaway struct {
	ID                   uuid.UUID
	ChannelID            string
	Type                 GiveawayType
	Keyword              *string
	MinWatchedTime       *int64
	MinMessages          *int32
	MinUsedChannelPoints *int64
	MinFollowDuration    *int64
	RequireSubscription  bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
	StartedAt            *time.Time
	StoppedAt            *time.Time
	CreatedByUserID      string

	isNil bool
}

func (g Giveaway) IsNil() bool {
	return g.isNil
}

var GiveawayNil = Giveaway{
	isNil: true,
}

type GiveawayWinner struct {
	UserID      string
	UserLogin   string
	DisplayName string

	isNil bool
}

func (w GiveawayWinner) IsNil() bool {
	return w.isNil
}

var ChannelGiveawayWinnerNil = GiveawayWinner{
	isNil: true,
}

type GiveawayParticipant struct {
	ID          uuid.UUID
	UserID      string
	UserLogin   string
	DisplayName string
	IsWinner    bool
	GiveawayID  uuid.UUID

	isNil bool
}
