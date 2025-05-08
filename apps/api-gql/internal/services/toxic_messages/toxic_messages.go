package toxic_messages

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	toxicmessagesrepository "github.com/twirapp/twir/libs/repositories/toxic_messages"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Repository toxicmessagesrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		repository: opts.Repository,
	}
}

type Service struct {
	repository toxicmessagesrepository.Repository
}

type GetListInput struct {
	Page    int
	PerPage int
}

type GetListOutput struct {
	Items []entity.ToxicMessage
	Total int
}

func (c *Service) GetList(ctx context.Context, input GetListInput) (GetListOutput, error) {
	data, err := c.repository.GetList(
		ctx,
		toxicmessagesrepository.GetListInput{
			Page:    input.Page,
			PerPage: input.PerPage,
		},
	)
	if err != nil {
		return GetListOutput{}, err
	}

	converted := make([]entity.ToxicMessage, 0, len(data.Items))
	for _, m := range data.Items {
		converted = append(
			converted, entity.ToxicMessage{
				ID:              m.ID,
				ChannelID:       m.ChannelID,
				ReplyToToUserID: m.ReplyToToUserID,
				Text:            m.Text,
				CreatedAt:       m.CreatedAt,
			},
		)
	}

	return GetListOutput{
		Items: converted,
		Total: data.Total,
	}, nil
}
