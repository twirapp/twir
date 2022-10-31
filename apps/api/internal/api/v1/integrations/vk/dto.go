package vk

type vkDataDto struct {
	UserId string `validate:"required" json:"userId"`
}

type vkDto struct {
	Enabled *bool     `validate:"required" json:"enabled"`
	Data    vkDataDto `validate:"required" json:"data"`
}
