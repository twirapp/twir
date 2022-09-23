package types

import model "tsuwari/models"

type Timer struct {
	Model *model.ChannelsTimers
	
	SendIndex int
}

type Store map[string]*Timer