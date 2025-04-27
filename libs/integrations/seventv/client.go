package seventv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/twirapp/twir/libs/integrations/seventv/api"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewClient(token string) Client {
	// client with header

	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			&headerTransport{
				base:  http.DefaultTransport,
				token: token,
			},
		),
	}

	client := graphql.NewClient("https://api.7tv.app/v4/gql", httpClient)

	return Client{
		client: client,
	}
}

type Client struct {
	client graphql.Client
}

func unwrapGqlErrorsList(err error) error {
	var listError gqlerror.List
	switch {
	case nil == err:
		return nil
	case errors.As(err, &listError):
		unwrapped := make([]string, len(listError))
		for i, e := range listError.Unwrap() {
			var castedErr *gqlerror.Error
			castOk := errors.As(e, &castedErr)
			if castOk {
				unwrapped[i] = castedErr.Message
			} else {
				unwrapped[i] = e.Error()
			}
		}

		return errors.New(strings.Join(unwrapped, ", "))
	default:
		return fmt.Errorf("%w: %s", err, err.Error())
	}
}

func (c *Client) SearchEmote(ctx context.Context, query string) (
	*api.SearchEmoteByNameResponse,
	error,
) {
	res, err := api.SearchEmoteByName(ctx, c.client, query)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetProfileByTwitchId(
	ctx context.Context,
	id string,
) (*api.GetProfileByTwitchIdResponse, error) {
	res, err := api.GetProfileByTwitchId(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetOneEmoteById(ctx context.Context, id string) (
	*api.GetOneEmoteByIdResponse,
	error,
) {
	res, err := api.GetOneEmoteById(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetOneEmoteByNameOrLink(ctx context.Context, query string) (
	api.TwirSeventvEmote,
	error,
) {
	if id := FindEmoteIdInInput(query); id != "" {
		resp, err := c.GetOneEmoteById(ctx, id)
		if err != nil {
			return api.TwirSeventvEmote{}, unwrapGqlErrorsList(err)
		}

		return resp.Emotes.Emote.TwirSeventvEmote, nil
	}

	resp, err := c.SearchEmote(ctx, query)
	if err != nil {
		return api.TwirSeventvEmote{}, unwrapGqlErrorsList(err)
	}

	if len(resp.Emotes.Search.Items) == 0 {
		return api.TwirSeventvEmote{}, fmt.Errorf("emote not found")
	}

	return resp.Emotes.Search.Items[0].TwirSeventvEmote, nil
}

func (c *Client) AddEmote(
	ctx context.Context,
	emoteSetID string,
	emoteID string,
	alias string,
) error {
	_, err := api.AddEmoteToSet(
		ctx,
		c.client,
		emoteSetID,
		api.EmoteSetEmoteId{
			EmoteId: emoteID,
			Alias:   &alias,
		},
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) RemoveEmote(ctx context.Context, emoteSetId, emoteName, emoteId string) error {
	_, err := api.DeleteEmoteFromSet(
		ctx,
		c.client,
		emoteSetId,
		api.EmoteSetEmoteId{EmoteId: emoteId, Alias: &emoteName},
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) RenameEmote(ctx context.Context, emoteSetId, emoteId, alias string) error {
	_, err := api.RenameEmote(
		ctx,
		c.client,
		emoteSetId,
		api.EmoteSetEmoteId{EmoteId: emoteId},
		alias,
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) GetProfileById(ctx context.Context, id string) (
	*api.GetProfileByIdResponse,
	error,
) {
	res, err := api.GetProfileById(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}
