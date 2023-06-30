package giveaways

type postGiveawayDto struct {
	Running               bool     `validate:"required" json:"running"`
	GiveawayType          string   `validate:"required" json:"giveawayType"`
	Keyword               *string  `json:"keyword,omitempty"`
	IsNeedAnnounce        bool     `validate:"required" json:"isNeedAnnounce"`
	MinimumWatchTime      int      `validate:"required,min:0" json:"minimumWatchTime"`
	MinimumFollowTime     int      `validate:"required,min:0" json:"minimumFollowTime"`
	MinimumMessages       int      `validate:"required,min:0" json:"minimumMessages"`
	MinimumSubTier        int      `validate:"required,min:0" json:"minimumSubTier"`
	MinimumSubTime        int      `validate:"required,min:0" json:"minimumSubTime"`
	RolesIds              []string `validate:"required" json:"rolesIds"`
	GiveawayMinimumNumber *int     `json:"giveawayMinimumNumber,omitempty"`
	GiveawayMaximumNumber *int     `json:"giveawayMaximumNumber,omitempty"`
	WinnersCount          int      `validate:"required,min:0" json:"winnersCount"`
	SubLuck               int      `validate:"required,min:0,max:10" json:"subLuck"`
	SubTier1Luck          int      `validate:"required,min:0,max:10" json:"subTier1Luck"`
	SubTier2Luck          int      `validate:"required,min:0,max:10" json:"subTier2Luck"`
	SubTier3Luck          int      `validate:"required,min:0,max:10" json:"subTier3Luck"`
}
