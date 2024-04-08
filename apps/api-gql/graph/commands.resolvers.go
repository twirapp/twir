package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.45

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	model "github.com/satont/twir/libs/gomodels"
	gmodel "github.com/twirapp/twir/apps/api-gql/graph/model"
)

// CreateCommand is the resolver for the createCommand field.
func (r *mutationResolver) CreateCommand(ctx context.Context, name string, description *string) (*gmodel.Command, error) {
	newCommand := model.ChannelsCommands{
		ID:          uuid.NewString(),
		Name:        name,
		Description: null.StringFromPtr(description),
	}

	commands = append(
		commands, newCommand,
	)

	return &gmodel.Command{
		ID:   newCommand.ID,
		Name: newCommand.Name,
	}, nil
}

// Commands is the resolver for the commands field.
func (r *queryResolver) Commands(ctx context.Context) ([]*gmodel.Command, error) {
	responseCommands := []*gmodel.Command{}

	for _, command := range commands {
		responses := []*gmodel.CommandResponse{}
		for _, response := range command.Responses {
			responses = append(
				responses, &gmodel.CommandResponse{
					ID:        response.ID,
					CommandID: response.CommandID,
					Text:      response.Text.String,
					Order:     response.Order,
				},
			)
		}
		responseCommands = append(
			responseCommands, &gmodel.Command{
				ID:          command.ID,
				Name:        command.Name,
				Description: command.Description.Ptr(),
				Responses:   responses,
				Aliases:     command.Aliases,
			},
		)
	}

	return responseCommands, nil
}

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
var commands = []model.ChannelsCommands{}
