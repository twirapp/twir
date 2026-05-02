package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func GreetingEntityTo(e entity.Greeting) gqlmodel.Greeting {
	return gqlmodel.Greeting{
		ID:           e.ID,
		UserID:       e.UserID.String(),
		Enabled:      e.Enabled,
		IsReply:      e.IsReply,
		Text:         e.Text,
		WithShoutOut: e.WithShoutOut,
	}
}
