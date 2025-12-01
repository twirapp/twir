package model

type ChannelIntegrationLastfm struct {
	ID         int
	ChannelID  string
	Enabled    bool
	SessionKey *string
	UserName   *string
	Avatar     *string

	isNil bool
}

func (c ChannelIntegrationLastfm) IsNil() bool {
	return c.isNil
}

var Nil = ChannelIntegrationLastfm{
	isNil: true,
}
