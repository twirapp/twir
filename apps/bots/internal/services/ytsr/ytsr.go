package ytsr

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/goccy/go-json"
	"github.com/raitonoberu/ytsearch"
	"github.com/samber/lo"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/ytsr"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Config  cfg.Config
	Logger  *slog.Logger
	TwirBus *buscore.Bus
}

func New(opts Opts) error {
	s := &Service{
		ytRegexp: *regexp.MustCompile(
			`(?m)http(?:s?):\/\/(?:www\.)?youtu(?:be\.com\/watch\?v=|\.be\/)([\w\-\_]*)(&(amp;)?‌​[\w\?‌​=]*)?`,
		),
		config: opts.Config,
		logger: opts.Logger,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				opts.Logger.Info("ytsr started")
				return opts.TwirBus.YTSRSearch.SubscribeGroup("ytsr", s.search)
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.YTSRSearch.Unsubscribe()
				return nil
			},
		},
	)

	return nil
}

type Service struct {
	ytRegexp    regexp.Regexp
	linksRegexp regexp.Regexp

	config cfg.Config
	logger *slog.Logger
}

type internalSong struct {
	odesliUrl    string
	youtubeQuery string
}

func (c *Service) search(ctx context.Context, req ytsr.SearchRequest) (
	ytsr.SearchResponse,
	error,
) {
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
					c.logger.Error("searchOdesli", logger.Error(err))
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
	} else if !req.OnlyLinks {
		internalSongs = append(
			internalSongs,
			internalSong{
				youtubeQuery: req.Search,
			},
		)
	}

	if len(internalSongs) == 0 {
		return ytsr.SearchResponse{}, nil
	}

	var wg sync.WaitGroup
	var songsMu sync.Mutex
	songs := make([]ytsr.Song, 0, len(internalSongs))

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
				c.logger.Error("searchByText", logger.Error(err))
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
				ytsr.Song{
					Title:        res.Title,
					Id:           res.ID,
					Views:        uint64(res.ViewCount),
					Duration:     uint64(res.Duration),
					ThumbnailUrl: videoThumbNail,
					IsLive:       false,
					Author: ytsr.SongAuthor{
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

	return ytsr.SearchResponse{
		Songs: songs,
	}, nil
}

func (c *Service) searchByText(_ context.Context, query string) (ytsearch.VideoItem, error) {
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

type odesliPlatform struct {
	Url string `json:"url"`
}
type odesliPlatforms map[string]*odesliPlatform

type odesliResponse struct {
	PageUrl         string          `json:"pageUrl"`
	LinksByPlatform odesliPlatforms `json:"linksByPlatform"`
}

type odesliErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Code       string `json:"code"`
}

func (c *Service) searchOdesli(ctx context.Context, query string) (odesliResponse, error) {
	result := odesliResponse{}

	reqUrl := fmt.Sprintf(
		"https://api.song.link/v1-alpha.1/links?url=%s&key=%s",
		query,
		c.config.OdesliApiKey,
	)
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return result, err
	}
	req = req.WithContext(ctx)
	resp, err := httpClient.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != 200 {
		odesliError := odesliErrorResponse{}
		if err = json.Unmarshal(body, &odesliError); err != nil {
			return result, err
		}
		return result, fmt.Errorf(`odesli error for input "%s": %s`, query, odesliError.Code)
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	if _, ok := result.LinksByPlatform["youtube"]; !ok {
		return result, fmt.Errorf(`odesli error for input "%s": %s`, query, "no youtube link")
	}

	return result, nil
}
