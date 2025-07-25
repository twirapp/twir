package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"

	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/services/toxic_messages"
)

// AdminToxicMessages is the resolver for the adminToxicMessages field.
func (r *queryResolver) AdminToxicMessages(ctx context.Context, input gqlmodel.AdminToxicMessagesInput) (*gqlmodel.AdminToxicMessagesPayload, error) {
	var (
		page    int
		perPage int
	)

	if input.Page.IsSet() {
		page = *input.Page.Value()
	}
	if input.PerPage.IsSet() {
		perPage = *input.PerPage.Value()
	}

	list, err := r.deps.ToxicMessagesService.GetList(
		ctx,
		toxic_messages.GetListInput{
			Page:    page,
			PerPage: perPage,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error getting toxic messages: %w", err)
	}

	converted := make([]gqlmodel.ToxicMessage, 0, len(list.Items))
	for _, m := range list.Items {
		converted = append(
			converted,
			gqlmodel.ToxicMessage{
				ID:            m.ID,
				ChannelID:     m.ChannelID,
				ReplyToUserID: m.ReplyToToUserID,
				Text:          m.Text,
				CreatedAt:     m.CreatedAt,
			},
		)
	}

	return &gqlmodel.AdminToxicMessagesPayload{
		Items: converted,
		Total: list.Total,
	}, nil
}

// ChannelProfile is the resolver for the channelProfile field.
func (r *toxicMessageResolver) ChannelProfile(ctx context.Context, obj *gqlmodel.ToxicMessage) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.ChannelID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.ChannelID)
}

// ToxicMessage returns graph.ToxicMessageResolver implementation.
func (r *Resolver) ToxicMessage() graph.ToxicMessageResolver { return &toxicMessageResolver{r} }

type toxicMessageResolver struct{ *Resolver }
