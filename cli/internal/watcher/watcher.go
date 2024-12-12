package watcher

import (
	"strings"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/samber/lo"
)

type Watcher struct {
	chann chan notify.EventInfo
}

func New() *Watcher {
	return &Watcher{
		chann: make(chan notify.EventInfo, 1),
	}
}

func (c *Watcher) Start(path string) (chan struct{}, error) {
	watchPath := path + "..."

	updateChannel := make(chan struct{}, 1)

	if err := notify.Watch(watchPath, c.chann, notify.All); err != nil {
		return nil, err
	}

	reload, _ := lo.NewDebounce(
		1*time.Second,
		func() {
			updateChannel <- struct{}{}
		},
	)

	go func() {
		for event := range c.chann {
			if strings.HasSuffix(event.Path(), "~") || strings.Contains(event.Path(), ".out") {
				continue
			}

			reload()
		}
	}()

	return updateChannel, nil
}

func (c *Watcher) Stop() {
	notify.Stop(c.chann)
}
