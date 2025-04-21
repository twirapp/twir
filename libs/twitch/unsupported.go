package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

const categorySearchQuery = `[{"operationName":"EditBroadcastCategoryDropdownSearch","variables":{"query":"%s"},"extensions":{"persistedQuery":{"version":1,"sha256Hash":"ccad6fa3d84008d710f2d69f7f9bcbd30d6b54ed1cecea5fd9a0a28e3f2508c7"}}}]`

type TwitchGqlSearchCategoryResponse struct {
	Data *struct {
		SearchCategories *struct {
			Edges []struct {
				Cursor string `json:"cursor"`
				Node   struct {
					BoxArtURL                           string    `json:"boxArtURL"`
					DisplayName                         string    `json:"displayName"`
					Id                                  string    `json:"id"`
					Name                                string    `json:"name"`
					ViewersCount                        int       `json:"viewersCount"`
					FollowersCount                      int       `json:"followersCount"`
					IsRestrictedForCurrentUserAndRegion bool      `json:"isRestrictedForCurrentUserAndRegion"`
					IsMature                            bool      `json:"isMature"`
					OriginalReleaseDate                 time.Time `json:"originalReleaseDate"`
					Platforms                           []string  `json:"platforms"`
					Publishers                          []string  `json:"publishers"`
					Tags                                []struct {
						Id       string `json:"id"`
						Typename string `json:"__typename"`
					} `json:"tags"`
					Typename string `json:"__typename"`
				} `json:"node"`
				Typename string `json:"__typename"`
			} `json:"edges"`
			Typename string `json:"__typename"`
		} `json:"searchCategories"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
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
		searchResponse[0].Data.SearchCategories == nil ||
		len(searchResponse[0].Data.SearchCategories.Edges) == 0 {
		return nil, fmt.Errorf("cannot get game from twitch: empty response")
	}

	return &FoundCategory{
		ID:   searchResponse[0].Data.SearchCategories.Edges[0].Node.Id,
		Name: searchResponse[0].Data.SearchCategories.Edges[0].Node.Name,
	}, nil
}
