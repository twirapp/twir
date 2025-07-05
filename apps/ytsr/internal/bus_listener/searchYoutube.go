package bus_listener

import (
	"context"

	"github.com/raitonoberu/ytsearch"
)

func (c *YtsrServer) searchByText(_ context.Context, query string) (ytsearch.VideoItem, error) {
	q := ytsearch.VideoSearch(query)

	items, err := q.Next()
	if err != nil {
		return ytsearch.VideoItem{}, nil
	}

	if len(items.Videos) == 0 {
		return ytsearch.VideoItem{}, nil
	}

	return *items.Videos[0], nil
}
