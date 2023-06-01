package twitch

import (
	"errors"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"sync"
)

func (c *helixClient) getManyUsers(ids, names []string) ([]helix.User, error) {
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	chunks := lo.Chunk(lo.If(len(ids) > 0, ids).Else(names), 100)
	wg.Add(len(chunks))

	var users []helix.User
	var usersMutex sync.Mutex

	for _, chunk := range chunks {
		go func(chunk []string) {
			defer wg.Done()

			param := lo.If(len(ids) > 0, &helix.UsersParams{IDs: chunk}).Else(&helix.UsersParams{Logins: names})
			req, err := c.GetUsers(param)

			if err != nil || req.ErrorMessage != "" {
				errCh <- lo.If(err != nil, err).Else(errors.New(req.ErrorMessage))
				return
			}

			usersMutex.Lock()
			users = append(users, req.Data.Users...)
			usersMutex.Unlock()

			errCh <- nil
		}(chunk)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for range chunks {
		if err := <-errCh; err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (c *helixClient) GetManyUsersByIds(ids []string) ([]helix.User, error) {
	return c.getManyUsers(ids, nil)
}

func (c *helixClient) GetManyUsersByNames(names []string) ([]helix.User, error) {
	return c.getManyUsers(nil, names)
}
