package app

import (
	cfg "github.com/twirapp/twir/libs/config"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
)

func newVKVideoChatClient(config cfg.Config) (*vkintegrations.VideoChatClient, error) {
	return vkintegrations.NewVideoChatClient(vkintegrations.VideoChatClientOpts{APIBaseURL: config.VKVideoDevAPIBaseURL})
}
