package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

// TwirStats is the resolver for the twirStats field.
func (r *queryResolver) TwirStats(ctx context.Context) (*gqlmodel.TwirStats, error) {
	return r.twirStats.GetCachedData(), nil
}