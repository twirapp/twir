package users

import (
	"net/http"
	"strings"
	"sync"

	"github.com/satont/tsuwari/libs/twitch"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/api/internal/types"
)

type RequestUser struct {
	ID   *string
	Name *string
}

func handleGet(ids string, names string, services *types.Services) ([]helix.User, error) {
	twitchClient, err := twitch.NewAppClient(*services.Config, services.Grpc.Tokens)
	if err != nil {
		services.Logger.Error(err)
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
		services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get users")
	default:
		return users, nil
	}
}
