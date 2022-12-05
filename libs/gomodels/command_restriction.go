package model

type CommandRestriction struct {
	ID        string `gorm:"type:uuid;column:id"        json:"id"`
	CommandID string `gorm:"type:text;column:commandId" json:"commandId"`
	Type      string `gorm:"type:text;column:type"      json:"type"`
	Value     string `gorm:"type:text;column:value"     json:"value"`
}

func (CommandRestriction) TableName() string {
	return "commands_restrictions"
}
