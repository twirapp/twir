package model

type ChannelIntegrationValorant struct {
	ID                   int
	ChannelID            string
	Enabled              bool
	AccessToken          *string
	RefreshToken         *string
	UserName             *string
	ValorantActiveRegion *string
	ValorantPuuid        *string

	isNil bool
}

func (c ChannelIntegrationValorant) IsNil() bool {
	return c.isNil
}

var Nil = ChannelIntegrationValorant{
	isNil: true,
}
