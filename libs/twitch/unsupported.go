package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/imroc/req/v3"
)

const twitchOfficialClientId = "kimne78kx3ncx6brgo4mv6wki5h1ko"

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
		SetHeader("client-id", twitchOfficialClientId).
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

const firstFollowerQuery = `[
	{
		"extensions": {
			"persistedQuery": {
				"sha256Hash": "3316194bb52051e2f9184012f6171b9aed4d457994568f1b4ed4a11e37a18b5c",
				"version": 1
			}
		},
		"operationName": "Followers",
		"variables": {
			"limit": 50,
			"login": "%s",
			"order": "ASC"
		}
	}
]`

type FollowersResponse struct {
	Data *struct {
		User *struct {
			Id        string `json:"id"`
			Followers *struct {
				Edges []struct {
					Cursor string `json:"cursor"`
					Node   *struct {
						Id              string `json:"id"`
						BannerImageURL  string `json:"bannerImageURL"`
						DisplayName     string `json:"displayName"`
						Login           string `json:"login"`
						ProfileImageURL string `json:"profileImageURL"`
						Self            struct {
							CanFollow bool        `json:"canFollow"`
							Follower  interface{} `json:"follower"`
							Typename  string      `json:"__typename"`
						} `json:"self"`
						Typename string `json:"__typename"`
					} `json:"node"`
					Typename string `json:"__typename"`
				} `json:"edges"`
				PageInfo struct {
					HasNextPage bool   `json:"hasNextPage"`
					Typename    string `json:"__typename"`
				} `json:"pageInfo"`
				Typename string `json:"__typename"`
			} `json:"followers"`
			Typename string `json:"__typename"`
		} `json:"user"`
	} `json:"data"`
	Extensions struct {
		DurationMilliseconds int    `json:"durationMilliseconds"`
		OperationName        string `json:"operationName"`
		RequestID            string `json:"requestID"`
	} `json:"extensions"`
}

type FirstFollower struct {
	Id          string
	Login       string
	DisplayName string
}

func GetFirstChannelFollowers(ctx context.Context, channelName string) ([]FirstFollower, error) {
	var data []FollowersResponse

	res, err := req.
		SetContext(ctx).
		SetHeader("client-id", twitchOfficialClientId).
		SetBodyJsonString(fmt.Sprintf(firstFollowerQuery, channelName)).
		SetSuccessResult(&data).
		Post("https://gql.twitch.tv/gql")
	if err != nil {
		return nil, err
	}
	if !res.IsSuccessState() {
		return nil, fmt.Errorf("cannot get followers from twitch: %s", res.String())
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("cannot get followers from twitch: empty response")
	}

	firstEdge := data[0]

	if firstEdge.Data == nil || firstEdge.Data.User == nil ||
		firstEdge.Data.User.Followers == nil || len(firstEdge.Data.User.Followers.Edges) == 0 {
		return nil, fmt.Errorf("cannot get followers from twitch: empty response")
	}

	followers := make([]FirstFollower, 0, len(firstEdge.Data.User.Followers.Edges))
	for _, edge := range firstEdge.Data.User.Followers.Edges {
		if edge.Node == nil {
			continue
		}

		followers = append(
			followers, FirstFollower{
				Id:          edge.Node.Id,
				Login:       edge.Node.Login,
				DisplayName: edge.Node.DisplayName,
			},
		)
	}

	return followers, nil
}
