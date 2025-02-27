package song

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

var Playlist = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "song playlist",
		Description: null.StringFrom("Shows current song playlist (if any)"),
		RolesIDS:    pq.StringArray{},
		Module:      "SONGS",
		Visible:     true,
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		song := parseCtx.Cacher.GetCurrentSong(ctx)
		if song == nil || song.Playlist == nil {
			return nil, nil
		}

		info := ""
		if song.Playlist.Name != nil {
			info = fmt.Sprintf(`%s (%d followers)`, *song.Playlist.Name, *song.Playlist.Followers)
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"%s %s",
					info,
					song.Playlist.Href,
				),
			},
		}

		return result, nil
	},
}
