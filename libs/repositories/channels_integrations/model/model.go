package model

type ChannelIntegration struct {
	ID            string
	ChannelID     string
	IntegrationID string
	Enabled       bool
	AccessToken   *string
	RefreshToken  *string
	Data          *Data

	isNil bool
}

func (c ChannelIntegration) IsNil() bool {
	return c.isNil
}

type Data struct {
	UserName *string `json:"userName"`
	Avatar   *string `json:"avatar"`
}

var Nil = ChannelIntegration{
	isNil: true,
}
