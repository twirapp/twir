package resolvers

import (
	"github.com/twirapp/twir/apps/api-gql/gqlmodel"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	NewCommandChann chan *gqlmodel.Command
}
