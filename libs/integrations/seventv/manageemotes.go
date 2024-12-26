package seventv

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/imroc/req/v3"
)

const addMutation = `
mutation AddEmoteToSet($setId: Id!, $emote: EmoteSetEmoteId!) {
	emoteSets {
		emoteSet(id: $setId) {
			addEmote(id: $emote) {
				id
				__typename
			}
		__typename
		}
	__typename
	}
}`
const removeMutation = `
mutation RemoveEmoteFromSet($setId: Id!, $emote: EmoteSetEmoteId!) {
	emoteSets {
		emoteSet(id: $setId) {
			removeEmote(id: $emote) {
				id
				__typename
			}
		__typename
	}
		__typename
	}
}
`

type requestBody struct {
	OperationName string               `json:"operationName"`
	Query         string               `json:"query"`
	Variables     requestBodyVariables `json:"variables"`
}
type requestBodyVariables struct {
	Emote map[string]any `json:"emote"`
	SetId string         `json:"setId"`
}

var ErrCannotModify = errors.New("cannot modify channel emote set")
var ErrBadEmoteUrl = errors.New("bad emote url")
var ErrCannotAdd = errors.New("cannot add emote")

var emoteRegex = regexp.MustCompile(`((cdn.)?7tv.app/emotes/)(?P<id>.{26})`)

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

type sevenTvError struct {
	Message    string `json:"message"`
	Extensions map[string]any
}

type sevenTvResponse struct {
	Errors []sevenTvError `json:"errors"`
}

func emoteAction(ctx context.Context, action, sevenTvToken, input, setID string) error {
	emoteId := FindEmoteIdInInput(input)
	if emoteId == "" {
		return ErrBadEmoteUrl
	}

	var operationName string
	if action == "ADD" {
		operationName = "AddEmoteToSet"
	} else {
		operationName = "RemoveEmoteFromSet"
	}

	var query string
	if action == "ADD" {
		query = addMutation
	} else {
		query = removeMutation
	}

	body := requestBody{
		OperationName: operationName,
		Query:         query,
		Variables: requestBodyVariables{
			Emote: nil,
			SetId: setID,
		},
	}

	if action == "ADD" {
		body.Variables.Emote = map[string]any{
			"emoteId": emoteId,
		}
	} else {
		body.Variables.Emote = map[string]any{
			"emoteId": emoteId,
		}
	}

	var result sevenTvResponse
	resp, err := req.
		SetContext(ctx).
		SetBody(body).
		SetBearerAuthToken(sevenTvToken).
		SetSuccessResult(&result).
		Post("https://7tv.io/v4/gql")
	if err != nil {
		return nil
	}
	if !resp.IsSuccessState() || len(result.Errors) > 0 {
		if len(result.Errors) > 0 {
			var errs []string
			for _, err := range result.Errors {
				errs = append(errs, err.Message)
			}

			return fmt.Errorf("%w: %s", ErrCannotAdd, errs)
		}

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

func RenameEmote(ctx context.Context, sevenTvToken, setID, emoteID, newName string) error {
	// body := map[string]any{
	// 	"operationName": "ChangeEmoteInSet",
	// 	"variables": map[string]string{
	// 		"action":   "UPDATE",
	// 		"id":       setID,
	// 		"emote_id": emoteID,
	// 		"name":     newName,
	// 	},
	// 	"query": query,
	// }
	//
	// var result sevenTvResponse
	// resp, err := req.
	// 	SetContext(ctx).
	// 	SetBody(body).
	// 	SetBearerAuthToken(sevenTvToken).
	// 	SetSuccessResult(&result).
	// 	Post("https://7tv.io/v3/gql")
	// if err != nil {
	// 	return nil
	// }
	// if !resp.IsSuccessState() || len(result.Errors) > 0 {
	// 	return fmt.Errorf("%w: %s", ErrCannotModify, resp.String())
	// }
	//
	// return nil
	return nil
}

func RemoveEmoteByID(ctx context.Context, sevenTvToken, setID, emoteID string) error {
	// body := map[string]any{
	// 	"operationName": "ChangeEmoteInSet",
	// 	"variables": map[string]string{
	// 		"action":   "REMOVE",
	// 		"id":       setID,
	// 		"emote_id": emoteID,
	// 	},
	// 	"query": query,
	// }
	//
	// var result sevenTvResponse
	// resp, err := req.
	// 	SetContext(ctx).
	// 	SetBody(body).
	// 	SetBearerAuthToken(sevenTvToken).
	// 	SetSuccessResult(&result).
	// 	Post("https://7tv.io/v3/gql")
	// if err != nil {
	// 	return nil
	// }
	// if !resp.IsSuccessState() || len(result.Errors) > 0 {
	// 	return fmt.Errorf("%w: %s", ErrCannotModify, resp.String())
	// }
	//
	// return nil
	return nil
}
