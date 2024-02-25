package overlays

import (
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type DudesGrowRequest struct {
	websockets.DudesGrowRequest
}

type DudesUserSettings struct {
	DudeColor  *string      `json:"dudeColor"`
	DudeSprite *DudesSprite `json:"dudeSprite"`
	UserID     string       `json:"userId"`
	UserName   string       `json:"userName"`
	UserLogin  string       `json:"userLogin"`
}

type DudesSprite string

const (
	DudeSpriteRandom DudesSprite = "random"
	DudeSpriteAgent  DudesSprite = "agent"
	DudeSpriteCat    DudesSprite = "cat"
	DudeSpriteDude   DudesSprite = "dude"
	DudeSpriteGirl   DudesSprite = "girl"
	DudeSpriteSanta  DudesSprite = "santa"
	DudeSpriteSith   DudesSprite = "sith"
)

var AllDudesSpriteEnumValues = []DudesSprite{
	DudeSpriteAgent,
	DudeSpriteCat,
	DudeSpriteDude,
	DudeSpriteGirl,
	DudeSpriteSanta,
	DudeSpriteSith,
	DudeSpriteRandom,
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
	case DudeSpriteGirl:
		return DudeSpriteGirl.String()
	case DudeSpriteSanta:
		return DudeSpriteSanta.String()
	case DudeSpriteSith:
		return DudeSpriteSith.String()
	case DudeSpriteRandom:
		return DudeSpriteRandom.String()
	default:
		return ""
	}
}

func (c DudesSprite) IsValid() bool {
	switch c {
	case DudeSpriteAgent, DudeSpriteCat, DudeSpriteDude, DudeSpriteGirl, DudeSpriteSanta,
		DudeSpriteSith, DudeSpriteRandom:
		return true
	default:
		return false
	}
}
