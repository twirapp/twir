package seventv

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/imroc/req/v3"
)

const query = `mutation ChangeEmoteInSet($id: ObjectID!, $action: ListItemAction!, $emote_id: ObjectID!, $name: String) {
	emoteSet(id: $id) {
		emotes(id: $emote_id, action: $action, name: $name) {
			id
			name
		}
	}
}`

var ErrCannotModify = errors.New("cannot modify channel emote set")
var ErrBadEmoteUrl = errors.New("bad emote url")

var emoteRegex = regexp.MustCompile(`((cdn.)?7tv.app/emotes/)(?P<id>.{24})`)

func FindEmoteIdInInput(input string) string {
	var result string
	groupNames := emoteRegex.SubexpNames()
	for _, match := range emoteRegex.FindAllStringSubmatch(input, -1) {
		for groupIdx, group := range match {
			name := groupNames[groupIdx]
			if name == "id" {
				result = group
				break
			}
		}
	}

	return result
}

type sevenTvResponse struct {
	Errors []any `json:"errors"`
}

func emoteAction(ctx context.Context, action, sevenTvToken, input, setID string) error {
	emoteId := FindEmoteIdInInput(input)
	if emoteId == "" {
		return ErrBadEmoteUrl
	}

	body := map[string]any{
		"operationName": "ChangeEmoteInSet",
		"variables": map[string]string{
			"action":   action,
			"id":       setID,
			"emote_id": emoteId,
		},
		"query": query,
	}

	var result sevenTvResponse
	resp, err := req.
		SetContext(ctx).
		SetBody(body).
		SetBearerAuthToken(sevenTvToken).
		SetSuccessResult(&result).
		Post("https://7tv.io/v3/gql")
	if err != nil {
		return nil
	}
	if !resp.IsSuccessState() || len(result.Errors) > 0 {
		return fmt.Errorf("%w: %s", ErrCannotModify, resp.String())
	}

	return nil
}

func AddEmote(ctx context.Context, sevenTvToken, input, setID string) error {
	return emoteAction(ctx, "ADD", sevenTvToken, input, setID)
}

func RemoveEmote(ctx context.Context, sevenTvToken, input, setID string) error {
	return emoteAction(ctx, "REMOVE", sevenTvToken, input, setID)
}
