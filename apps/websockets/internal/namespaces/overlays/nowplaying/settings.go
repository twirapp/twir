package nowplaying

import (
	"github.com/satont/twir/apps/websockets/internal/protoutils"
)

func (c *NowPlaying) SendSettings(userId string) error {
	d, err := protoutils.CreateJsonWithProto(
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	return c.SendEvent(
		userId,
		"settings",
		d,
	)
}
