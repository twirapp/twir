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

func (c *Watcher) StartWithPaths(paths []string) (chan string, error) {
	updateChannel := make(chan string, 1)

	for _, path := range paths {
		watchPath := path + "..."
		if err := notify.Watch(watchPath, c.chann, notify.All); err != nil {
			return nil, err
		}
	}

	reload, _ := lo.NewDebounce(
		1*time.Second,
		func() {
			updateChannel <- ""
		},
	)

	var lastEventPath string

	go func() {
		for event := range c.chann {
			if strings.HasSuffix(event.Path(), "~") || strings.Contains(event.Path(), ".out") {
				continue
			}

			lastEventPath = event.Path()
			reload()
		}
	}()

	go func() {
		for range updateChannel {
			if lastEventPath != "" {
				updateChannel <- lastEventPath
				lastEventPath = ""
			}
		}
	}()

	return updateChannel, nil
}

func (c *Watcher) Stop() {
	notify.Stop(c.chann)
}

func (c *Watcher) WatchPaths(paths []string) error {
	for _, path := range paths {
		watchPath := path + "..."
		if err := notify.Watch(watchPath, c.chann, notify.All); err != nil {
			return err
		}
	}
	return nil
}

func (c *Watcher) GetEventPath() string {
	select {
	case event := <-c.chann:
		return event.Path()
	default:
		return ""
	}
}
