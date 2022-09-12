package model

type ChannelInegrationWithRelation struct {
	ChannelsIntegrations
	Integration Integrations `gorm:"foreignKey:IntegrationID"`
}
