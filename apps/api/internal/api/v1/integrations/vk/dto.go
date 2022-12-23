package vk

type vkDto struct {
	Code string `validate:"required" json:"code"`
}
