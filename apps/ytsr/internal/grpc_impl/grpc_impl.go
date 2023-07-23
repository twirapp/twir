package grpc_impl

import (
	"context"
	"fmt"
	"github.com/samber/lo"
	cfg "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/grpc/generated/ytsr"
	"go.uber.org/zap"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type YtsrServer struct {
	ytsr.UnimplementedYtsrServer

	ytRegexp    regexp.Regexp
	linksRegexp regexp.Regexp

	config cfg.Config
	logger *zap.Logger
}

func NewYtsrServer(config cfg.Config, logger *zap.Logger) *YtsrServer {
	return &YtsrServer{
		ytRegexp: *regexp.MustCompile(
			`(?m)http(?:s?):\/\/(?:www\.)?youtu(?:be\.com\/watch\?v=|\.be\/)([\w\-\_]*)(&(amp;)?‌​[\w\?‌​=]*)?`,
		),
		config: config,
		logger: logger,
	}
}

type internalSong struct {
	odesliUrl    string
	youtubeQuery string
}

func (c *YtsrServer) Search(ctx context.Context, req *ytsr.SearchRequest) (*ytsr.SearchResponse, error) {
	var linkMatches []string

	for _, part := range strings.Split(req.Search, " ") {
		u, err := url.Parse(part)
		if err != nil || u.Host == "" {
			continue
		}

		linkMatches = append(linkMatches, part)
	}

	var mu sync.Mutex
	internalSongs := make([]internalSong, 0, len(linkMatches))

	if len(linkMatches) > 0 {
		var wg sync.WaitGroup

		for _, link := range linkMatches {
			wg.Add(1)
			link := link
			go func() {
				defer wg.Done()

				odesliLink, err := c.searchOdesli(ctx, link)
				mu.Lock()
				defer mu.Unlock()

				// if odesli search fails, then we push raw youtube link to slice
				if err != nil {
					c.logger.Error("searchOdesli", zap.Error(err))
					internalSongs = append(
						internalSongs,
						internalSong{
							youtubeQuery: link,
						},
					)
					return
				}

				// push song with odesli link to slice

				internalSongs = append(
					internalSongs,
					internalSong{
						odesliUrl:    odesliLink.PageUrl,
						youtubeQuery: odesliLink.LinksByPlatform["youtube"].Url,
					},
				)
			}()
		}
		wg.Wait()
	} else {
		internalSongs = append(
			internalSongs,
			internalSong{
				youtubeQuery: req.Search,
			},
		)
	}

	var wg sync.WaitGroup
	var songsMu sync.Mutex
	songs := make([]*ytsr.Song, 0, len(internalSongs))

	for _, internalLink := range internalSongs {
		wg.Add(1)
		internalLink := internalLink

		youtubeMatch := c.ytRegexp.FindStringSubmatch(internalLink.youtubeQuery)

		go func() {
			defer wg.Done()

			res, err := c.searchByText(
				ctx,
				lo.IfF(
					len(youtubeMatch) != 0, func() string {
						return youtubeMatch[0]
					},
				).Else(internalLink.youtubeQuery),
			)
			if err != nil {
				c.logger.Error("searchByText", zap.Error(err))
				return
			}
			if res.ID == "" {
				return
			}

			videoThumbNail := lo.
				If[*string](len(res.Thumbnails) == 0, nil).
				Else(&res.Thumbnails[len(res.Thumbnails)-1].URL)
			channelThumbNail := lo.
				If[*string](len(res.Channel.Thumbnails) == 0, nil).
				Else(&res.Channel.Thumbnails[len(res.Channel.Thumbnails)-1].URL)

			songsMu.Lock()
			defer songsMu.Unlock()

			link := lo.
				If(internalLink.odesliUrl != "", internalLink.odesliUrl).
				Else(fmt.Sprintf("https://youtu.be/%s", res.ID))

			songs = append(
				songs,
				&ytsr.Song{
					Title:        res.Title,
					Id:           res.ID,
					Views:        uint64(res.ViewCount),
					Duration:     uint64(res.Duration),
					ThumbnailUrl: videoThumbNail,
					IsLive:       false,
					Author: &ytsr.SongAuthor{
						Name:      res.Channel.Title,
						ChannelId: res.Channel.ID,
						AvatarUrl: channelThumbNail,
					},
					Link: &link,
				},
			)
		}()
	}

	wg.Wait()

	return &ytsr.SearchResponse{
		Songs: songs,
	}, nil
}
