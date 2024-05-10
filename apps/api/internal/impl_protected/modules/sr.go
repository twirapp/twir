package modules

import (
	"context"
	"log/slog"
	"strings"
	"sync"

	"github.com/raitonoberu/ytsearch"
	"github.com/samber/lo"
	loParallel "github.com/samber/lo/parallel"
	"github.com/twirapp/twir/libs/api/messages/modules_sr"
)

func getThumbNailUrl(url string) string {
	return strings.Replace(url, "http://", "https://", 1)
}

func (c *Modules) ModulesSRSearchVideosOrChannels(
	_ context.Context,
	request *modules_sr.GetSearchRequest,
) (*modules_sr.GetSearchResponse, error) {
	response := &modules_sr.GetSearchResponse{
		Items: make([]*modules_sr.GetSearchResponse_Result, 0, len(request.Query)),
	}

	if len(request.Query) == 0 {
		return response, nil
	}

	var mu sync.Mutex

	loParallel.ForEach(
		request.Query, func(query string, _ int) {
			if query == "" {
				return
			}

			var search *ytsearch.SearchClient
			if request.Type == modules_sr.GetSearchRequest_CHANNEL {
				search = ytsearch.ChannelSearch(query)
			} else {
				search = ytsearch.VideoSearch(query)
			}

			res, err := search.Next()
			if err != nil {
				c.Logger.Error(
					"cannot find",
					slog.String("query", query),
				)
				return
			}

			mu.Lock()
			defer mu.Unlock()
			if request.Type == modules_sr.GetSearchRequest_CHANNEL {
				channels := lo.Map(
					res.Channels,
					func(item *ytsearch.ChannelItem, _ int) *modules_sr.GetSearchResponse_Result {
						thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)
						return &modules_sr.GetSearchResponse_Result{
							Id:        item.ID,
							Title:     item.Title,
							Thumbnail: thumb,
						}
					},
				)
				response.Items = append(
					response.Items,
					channels...,
				)
			} else {
				videos := lo.Map(
					res.Videos, func(item *ytsearch.VideoItem, _ int) *modules_sr.GetSearchResponse_Result {
						thumb := getThumbNailUrl(item.Thumbnails[len(item.Thumbnails)-1].URL)

						return &modules_sr.GetSearchResponse_Result{
							Id:        item.ID,
							Title:     item.Title,
							Thumbnail: thumb,
						}
					},
				)
				response.Items = append(
					response.Items,
					videos...,
				)
			}
		},
	)

	return response, nil
}
