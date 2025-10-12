package cache_invalidator

const (
	CacheInvalidatorSubject = "cache.invalidator"
)

type InvalidateRequest struct {
	Key string
}
