package internal

import (
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Services struct {
	DB     *gorm.DB
	Logger *zap.Logger
	Cfg    *cfg.Config

	BotsGrpc       bots.BotsClient
	TokensGrpc     tokens.TokensClient
	WebsocketsGrpc websockets.WebsocketClient
}

type Data struct {
	//
	UserName        string `json:"userName,omitempty"`
	UserDisplayName string `json:"userDisplayName,omitempty"`
	UserID          string `json:"userId,omitempty"`
	//
	RaidViewers int64 `json:"raidViewers,omitempty"`
	//
	ResubMonths  int64  `json:"resubMonths"`
	ResubStreak  int64  `json:"resubStreak"`
	ResubMessage string `json:"resubMessage"`
	SubLevel     string `json:"subLevel"`
	//
	OldStreamTitle    string `json:"oldStreamTitle"`
	NewStreamTitle    string `json:"newStreamTitle"`
	OldStreamCategory string `json:"oldStreamCategory"`
	NewStreamCategory string `json:"newStreamCategory"`
	//
	StreamTitle    string `json:"streamTitle"`
	StreamCategory string `json:"streamCategory"`
	//
	RewardID    string  `json:"-"'`
	RewardName  string  `json:"rewardName"`
	RewardCost  string  `json:"rewardCost"`
	RewardInput *string `json:"rewardInput"`
	//
	CommandName  string `json:"commandName"`
	CommandID    string `json:"-"`
	CommandInput string `json:"commandInput"`
	//
	TargetUserName        string `json:"targetUserName"`
	TargetUserDisplayName string `json:"targetUserDisplayName"`
	//
	DonateAmount   string `json:"donateAmount"`
	DonateMessage  string `json:"donateMessage"`
	DonateCurrency string `json:"donateCurrency"`
	//
	PrevOperation *DataFromPrevOperation `json:"prevOperation"`
	//
	KeywordID       string `json:"-"`
	KeywordName     string `json:"keywordName"`
	KeywordResponse string `json:"keywordResponse"`
	//
	GreetingText string `json:"greetingText"`
	//
	PollTitle                     string `json:"pollTitle"`
	PollOptionsNames              string `json:"pollOptionsNames"`
	PollTotalVotes                int    `json:"pollTotalVotes"`
	PollWinnerTitle               string `json:"pollWinnerTitle"`
	PollWinnerBitsVotes           int    `json:"pollWinnerBitsVotes"`
	PollWinnerChannelsPointsVotes int    `json:"pollWinnerChannelsPointsVotes"`
	PollWinnerTotalVotes          int    `json:"pollWinnerTotalVotes"`
	//
	PredictionTitle              string            `json:"predictionTitle"`
	PredictionOptionsNames       string            `json:"predictionOptionsNames"`
	PredictionTotalChannelPoints int               `json:"predictionTotalChannelPoints"`
	PredictionWinner             PredictionOutCome `json:"predictionWinner"`
	//

	ModeratorName        string `json:"moderatorName"`
	ModeratorDisplayName string `json:"moderatorDisplayName"`

	// ban
	BanReason        string `json:"banReason"`
	BanEndsInMinutes string `json:"banEndsInMinutes"`
}

type PredictionOutCome struct {
	Title       string `json:"title"`
	TotalUsers  int    `json:"totalUsers"`
	TotalPoints int    `json:"totalPoints"`
	TopUsers    string `json:"topUsers"`
}

type DataFromPrevOperation struct {
	UnmodedUserName string `json:"unmodedUserName"`
	UnvipedUserName string `json:"unvipedUserName"`
	BannedUserName  string `json:"bannedUserName"`
}
