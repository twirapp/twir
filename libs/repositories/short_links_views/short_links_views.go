package short_links_views

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, input CreateInput) error
	GetStatistics(ctx context.Context, input GetStatisticsInput) ([]StatisticsPoint, error)
	GetViews(ctx context.Context, input GetViewsInput) (GetViewsOutput, error)
	GetTopCountries(ctx context.Context, input GetTopCountriesInput) ([]CountryStats, error)
}

type CreateInput struct {
	ShortLinkID string
	UserID      *string
	IP          *string
	UserAgent   *string
	Country     *string
	City        *string
	CreatedAt   time.Time
}

type GetStatisticsInput struct {
	ShortLinkID string
	From        time.Time
	To          time.Time
	Interval    string // "hour" or "day"
}

type StatisticsPoint struct {
	Timestamp uint64
	Count     uint64
}

type GetViewsInput struct {
	ShortLinkID string
	Page        int
	PerPage     int
}

type View struct {
	ShortLinkID string
	UserID      *string
	Country     *string
	City        *string
	CreatedAt   time.Time
}

type GetViewsOutput struct {
	Views []View
	Total int
}

type GetTopCountriesInput struct {
	ShortLinkID string
	Limit       int
}

type CountryStats struct {
	Country string
	Count   uint64
}
