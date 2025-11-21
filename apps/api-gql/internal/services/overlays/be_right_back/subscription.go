package be_right_back

import "fmt"

func CreateSettingsSubscriptionKey(channelID string) string {
	return fmt.Sprintf("overlays:be_right_back:settings:%s", channelID)
}
