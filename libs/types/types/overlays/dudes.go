package overlays

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type DudesChangeColorRequest struct {
	websockets.DudesChangeColorRequest
}

type DudesGrowRequest struct {
	websockets.DudesGrowRequest
}

type DudesUserSettings struct {
	model.ChannelsOverlaysDudesUserSettings
}

type DudesSprite string

const (
	DudeSpriteAgent DudesSprite = "agent"
	DudeSpriteCat   DudesSprite = "cat"
	DudeSpriteDude  DudesSprite = "dude"
	DudeSprite      DudesSprite = "girl"
	DudeSpriteSanta DudesSprite = "santa"
	DudeSpriteSith  DudesSprite = "sith"
)

var AllDudesSpriteEnumValues = []DudesSprite{
	DudeSpriteAgent,
	DudeSpriteCat,
	DudeSpriteDude,
	DudeSprite,
	DudeSpriteSanta,
	DudeSpriteSith,
}

func (c DudesSprite) String() string {
	return string(c)
}

func (c DudesSprite) TSName() string {
	switch c {
	case DudeSpriteAgent:
		return DudeSpriteAgent.String()
	case DudeSpriteCat:
		return DudeSpriteCat.String()
	case DudeSpriteDude:
		return DudeSpriteDude.String()
	case DudeSprite:
		return DudeSprite.String()
	case DudeSpriteSanta:
		return DudeSpriteSanta.String()
	case DudeSpriteSith:
		return DudeSpriteSith.String()
	default:
		return ""
	}
}

func (c DudesSprite) IsValid() bool {
	switch c {
	case DudeSpriteAgent, DudeSpriteCat, DudeSpriteDude, DudeSprite, DudeSpriteSanta, DudeSpriteSith:
		return true
	default:
		return false
	}
}
