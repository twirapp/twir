package users

import (
	"net/http"
	"strings"
	"sync"

	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

func handleGetCategories(category string) ([]helix.Category, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	gameReq, err := twitchClient.GetGames(&helix.GamesParams{
		Names: []string{category},
	})
	if err != nil {
		return nil, err
	}

	categoriesResponse := make([]helix.Category, 0)

	if len(gameReq.Data.Games) != 0 {
		for _, game := range gameReq.Data.Games {
			categoriesResponse = append(categoriesResponse, helix.Category{
				ID:        game.ID,
				Name:      game.Name,
				BoxArtURL: game.BoxArtURL,
			})
		}

		return categoriesResponse, nil
	}

	if len(gameReq.Data.Games) == 0 {
		games, err := twitchClient.SearchCategories(&helix.SearchCategoriesParams{
			Query: category,
		})
		if err != nil {
			return nil, err
		}

		if len(games.Data.Categories) > 0 {
			return games.Data.Categories, nil
		}
	}

	return categoriesResponse, nil
}

type RequestUser struct {
	ID   *string
	Name *string
}

func handleGet(ids string, names string, services types.Services) ([]helix.User, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	config := do.MustInvoke[cfg.Config](di.Provider)

	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	usersForReq := []*RequestUser{}
	for _, v := range strings.Split(ids, ",") {
		if v == "" {
			continue
		}
		id := v
		usersForReq = append(usersForReq, &RequestUser{ID: &id})
	}
	for _, v := range strings.Split(names, ",") {
		if v == "" {
			continue
		}
		name := v
		usersForReq = append(usersForReq, &RequestUser{Name: &name})
	}

	if len(usersForReq) > 200 {
		return nil, fiber.NewError(400, "you cannot request more then 200 users")
	}

	if len(usersForReq) == 0 {
		return nil, nil
	}

	users := []helix.User{}

	chunks := lo.Chunk(usersForReq, 100)
	errCH := make(chan error)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	wg.Add(len(chunks))

	for _, c := range chunks {
		go func(c []*RequestUser) {
			defer wg.Done()

			usersByIds := lo.Filter(c, func(item *RequestUser, _ int) bool {
				return item.ID != nil
			})
			usersByNames := lo.Filter(c, func(item *RequestUser, _ int) bool {
				return item.Name != nil
			})

			req, err := twitchClient.GetUsers(&helix.UsersParams{
				IDs: lo.Map(usersByIds, func(item *RequestUser, _ int) string {
					return *item.ID
				}),
				Logins: lo.Map(usersByNames, func(item *RequestUser, _ int) string {
					return *item.Name
				}),
			})
			if err != nil {
				errCH <- err
			}
			mu.Lock()

			users = append(users, req.Data.Users...)
			mu.Unlock()
		}(c)
	}

	wg.Wait()

	select {
	case err := <-errCH:
		close(errCH)
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get users")
	default:
		return users, nil
	}
}
