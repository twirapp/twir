package twitch

import (
	"context"
	"fmt"

	"github.com/imroc/req/v3"
)

const categorySearchQuery = `[{"operationName":"SearchTray_SearchSuggestions",
"variables":{"queryFragment":"%s"},"extensions":{"persistedQuery":{"version":1,
"sha256Hash":"34e1899cd559b7d6a4ac25e3bdaad37a83324644b0085b4cc478d0f845f8f0de"}}}]`

type TwitchGqlSearchCategoryResponse struct {
	Data *struct {
		SearchSuggestions *struct {
			Edges []struct {
				Node *struct {
					Content *struct {
						Typename string `json:"__typename"`
						Game     *struct {
							ID   string `json:"id"`
							Slug string `json:"slug"`
						} `json:"game"`
					} `json:"content"`
					Text string `json:"text"`
				} `json:"node"`
			}
		}
	}
}

type FoundCategory struct {
	ID   string
	Name string
}

func SearchCategory(ctx context.Context, query string) (*FoundCategory, error) {
	var searchResponse []TwitchGqlSearchCategoryResponse
	res, err := req.
		SetContext(ctx).
		SetHeader("client-id", "kimne78kx3ncx6brgo4mv6wki5h1ko").
		SetSuccessResult(&searchResponse).
		SetBodyJsonString(fmt.Sprintf(categorySearchQuery, query)).
		Post("https://gql.twitch.tv/gql")
	if err != nil {
		return nil, err
	}
	if !res.IsSuccessState() {
		return nil, fmt.Errorf("cannot get game from twitch: %s", res.String())
	}

	if len(searchResponse) == 0 || searchResponse[0].Data == nil ||
		searchResponse[0].Data.SearchSuggestions == nil ||
		len(searchResponse[0].Data.SearchSuggestions.Edges) == 0 {
		return nil, fmt.Errorf("cannot get game from twitch: empty response")
	}

	var foundCategory *FoundCategory

	for _, edge := range searchResponse[0].Data.SearchSuggestions.Edges {
		if edge.Node == nil || edge.Node.Content == nil || edge.Node.Content.Game == nil {
			continue
		}

		if edge.Node.Content.Typename == "SearchSuggestionCategory" {
			foundCategory = &FoundCategory{
				ID:   edge.Node.Content.Game.ID,
				Name: edge.Node.Text,
			}
			break
		}
	}

	if foundCategory == nil {
		return nil, fmt.Errorf("cannot get game from twitch: category not found")
	}

	return foundCategory, nil
}
