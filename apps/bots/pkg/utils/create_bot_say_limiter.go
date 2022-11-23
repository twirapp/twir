package utils

import (
	"time"

	"github.com/aidenwallis/go-ratelimiting/local"
)

func CreateBotLimiter(isMod bool) local.SlidingWindow {
	if isMod {
		l, _ := local.NewSlidingWindow(20, 30*time.Second)
		return l
	} else {
		l, _ := local.NewSlidingWindow(1, 2*time.Second)
		return l
	}
}
