package roles

type roleSettingsDto struct {
	RequiredWatchTime         int64 `json:"requiredWatchTime"`
	RequiredMessages          int32 `json:"requiredMessages"`
	RequiredUsedChannelPoints int64 `json:"requiredUsedChannelPoints"`
}

type roleDto struct {
	Name        string          `validate:"required" json:"name"`
	Permissions []string        `json:"permissions"`
	Users       []string        `json:"users"`
	Settings    roleSettingsDto `json:"settings"`
}
