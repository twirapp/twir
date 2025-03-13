package entity

type SevenTvIntegrationData struct {
	IsEditor                   bool
	BotSeventvProfile          *SevenTvProfile
	UserSeventvProfile         *SevenTvProfile
	RewardIDForAddEmote        *string
	RewardIDForRemoveEmote     *string
	EmoteSetID                 *string
	DeleteEmotesOnlyAddedByApp bool
}

type SevenTvProfile struct {
	ID          string
	Username    string
	DisplayName string
	Editors     []SevenTvProfileEditor
	EmoteSetID  *string
	AvatarURI   string
}

type SevenTvProfileEditor struct {
	ID          string
	Permissions int
	Visible     bool
	AddedAt     int64
}
