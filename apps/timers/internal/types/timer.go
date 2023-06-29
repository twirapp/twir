package types

import model "github.com/satont/twir/libs/gomodels"

type Timer struct {
	Model *model.ChannelsTimers

	SendIndex int
}

type Store map[string]*Timer
