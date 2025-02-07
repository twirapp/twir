package entity

type TTSUserSettings struct {
	UserID         string
	Rate           int
	Pitch          int
	Volume         int
	Voice          string
	IsChannelOwner bool
}
