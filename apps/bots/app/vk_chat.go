package app

import vkintegrations "github.com/twirapp/twir/libs/integrations/vk"

func newVKVideoChatClient() (*vkintegrations.VideoChatClient, error) {
	return vkintegrations.NewVideoChatClient(vkintegrations.VideoChatClientOpts{})
}
