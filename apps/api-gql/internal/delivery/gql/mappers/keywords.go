package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func KeywordsFrom(k entity.Keyword) gqlmodel.Keyword {
	return gqlmodel.Keyword{
		ID:                  k.ID,
		Text:                k.Text,
		Response:            &k.Response,
		Enabled:             k.Enabled,
		Cooldown:            k.Cooldown,
		IsReply:             k.IsReply,
		IsRegularExpression: k.IsRegular,
		UsageCount:          k.Usages,
	}
}
