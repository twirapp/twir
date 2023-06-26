package youtube

import (
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/internal/namespaces/helpers"
	"github.com/satont/twir/apps/websockets/types"
)

type YouTube struct {
	manager  *melody.Melody
	services *types.Services
}

func NewYouTube(services *types.Services) *YouTube {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024 * 10
	youTube := &YouTube{
		manager:  m,
		services: services,
	}

	youTube.manager.HandleConnect(
		func(session *melody.Session) {
			err := helpers.CheckUserByApiKey(services.Gorm, session)
			if err != nil {
				services.Logger.Error(err)
			} else {
				youTube.handleConnect(session)
			}
		},
	)

	youTube.manager.HandleMessage(
		func(session *melody.Session, msg []byte) {
			youTube.handleMessage(session, msg)
		},
	)

	return youTube
}
