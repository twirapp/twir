package twitch

import (
	"errors"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"sync"
)

func (c *helixClient) GetManyStreams(channelIds []string) ([]helix.Stream, error) {
	errCh := make(chan error)
	wg := sync.WaitGroup{}

	chunks := lo.Chunk(channelIds, 100)
	wg.Add(len(chunks))

	var streams []helix.Stream
	var usersMutex sync.Mutex

	for _, chunk := range chunks {
		go func(chunk []string) {
			defer wg.Done()

			req, err := c.GetStreams(&helix.StreamsParams{
				UserIDs: chunk,
			})

			if err != nil || req.ErrorMessage != "" {
				errCh <- lo.If(err != nil, err).Else(errors.New(req.ErrorMessage))
				return
			}

			usersMutex.Lock()
			streams = append(streams, req.Data.Streams...)
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

	return streams, nil
}
