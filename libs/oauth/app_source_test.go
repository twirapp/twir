package oauth

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestAppTokenSourceCachesConcurrentFetch(t *testing.T) {
	var mu sync.Mutex
	calls := 0
	source, err := NewAppTokenSource(appFetch(func(context.Context, AppTokenKey) (AppToken, error) {
		mu.Lock()
		calls++
		mu.Unlock()
		return AppToken{AccessToken: "app", ObtainedAt: time.Unix(1, 0), ExpiresIn: time.Hour}, nil
	}), fixedClock{now: time.Unix(2, 0)}, 0)
	if err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	for range 8 {
		group.Add(1)
		go func() {
			defer group.Done()
			if _, err := source.Token(context.Background(), "kick-client"); err != nil {
				t.Error(err)
			}
		}()
	}
	group.Wait()
	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("fetches = %d", calls)
	}
}

type appFetch func(context.Context, AppTokenKey) (AppToken, error)

func (f appFetch) FetchAppToken(ctx context.Context, key AppTokenKey) (AppToken, error) {
	return f(ctx, key)
}
