package v3

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/twirapp/twir/libs/integrations/seventv/v3/api"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func NewClient(token string) Client {
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(
			&headerTransport{
				base:  http.DefaultTransport,
				token: token,
			},
		),
	}

	client := graphql.NewClient("https://7tv.io/v3/gql", httpClient)

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
	*api.SearchEmoteByNameV3Response,
	error,
) {
	res, err := api.SearchEmoteByNameV3(ctx, c.client, query)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetProfileByTwitchId(
	ctx context.Context,
	id string,
) (*api.GetProfileByTwitchIdV3Response, error) {
	res, err := api.GetProfileByTwitchIdV3(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetProfileByKickId(
	ctx context.Context,
	id string,
) (*api.GetProfileByKickIdV3Response, error) {
	res, err := api.GetProfileByKickIdV3(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetOneEmoteById(ctx context.Context, id string) (
	*api.GetOneEmoteByIdV3Response,
	error,
) {
	res, err := api.GetOneEmoteByIdV3(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

func (c *Client) GetOneEmoteByNameOrLink(ctx context.Context, query string) (
	api.TwirSeventvV3Emote,
	error,
) {
	if id := FindEmoteIdInInput(query); id != "" {
		resp, err := c.GetOneEmoteById(ctx, id)
		if err != nil {
			return api.TwirSeventvV3Emote{}, unwrapGqlErrorsList(err)
		}

		if resp.Emote == nil {
			return api.TwirSeventvV3Emote{}, fmt.Errorf("emote not found")
		}

		return resp.Emote.TwirSeventvV3Emote, nil
	}

	resp, err := c.SearchEmote(ctx, query)
	if err != nil {
		return api.TwirSeventvV3Emote{}, unwrapGqlErrorsList(err)
	}

	if len(resp.Emotes.Items) == 0 {
		return api.TwirSeventvV3Emote{}, fmt.Errorf("emote not found")
	}

	return resp.Emotes.Items[0].TwirSeventvV3Emote, nil
}

func (c *Client) AddEmote(
	ctx context.Context,
	emoteSetID string,
	emoteID string,
	alias string,
) error {
	name := &alias
	if alias == "" {
		name = nil
	}

	_, err := api.AddEmoteToSetV3(
		ctx,
		c.client,
		emoteSetID,
		emoteID,
		name,
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) RemoveEmote(ctx context.Context, emoteSetId, emoteId string) error {
	_, err := api.DeleteEmoteFromSetV3(
		ctx,
		c.client,
		emoteSetId,
		emoteId,
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) RenameEmote(ctx context.Context, emoteSetId, emoteId, alias string) error {
	_, err := api.RenameEmoteV3(
		ctx,
		c.client,
		emoteSetId,
		emoteId,
		alias,
	)

	return unwrapGqlErrorsList(err)
}

func (c *Client) GetProfileById(ctx context.Context, id string) (
	*api.GetProfileByIdV3Response,
	error,
) {
	res, err := api.GetProfileByIdV3(ctx, c.client, id)
	return res, unwrapGqlErrorsList(err)
}

var emoteIdRegex = regexp.MustCompile(`https?://7tv\.io/emotes/(\w+)`)

func FindEmoteIdInInput(input string) string {
	matches := emoteIdRegex.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
