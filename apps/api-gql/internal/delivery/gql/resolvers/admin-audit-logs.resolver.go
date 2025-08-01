package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"

	data_loader "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/dataloader"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/graph"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	audit_logs "github.com/twirapp/twir/apps/api-gql/internal/services/audit-logs"
)

// User is the resolver for the user field.
func (r *adminAuditLogResolver) User(ctx context.Context, obj *gqlmodel.AdminAuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.UserID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.UserID)
}

// Channel is the resolver for the channel field.
func (r *adminAuditLogResolver) Channel(ctx context.Context, obj *gqlmodel.AdminAuditLog) (*gqlmodel.TwirUserTwitchInfo, error) {
	if obj.ChannelID == nil {
		return nil, nil
	}

	return data_loader.GetHelixUserById(ctx, *obj.ChannelID)
}

// AdminAuditLogs is the resolver for the adminAuditLogs field.
func (r *queryResolver) AdminAuditLogs(ctx context.Context, input gqlmodel.AdminAuditLogsInput) (*gqlmodel.AdminAuditLogResponse, error) {
	var page int
	perPage := 20

	if input.Page.IsSet() {
		page = *input.Page.Value()
	}

	if input.PerPage.IsSet() {
		perPage = *input.PerPage.Value()
	}

	logsInput := audit_logs.GetManyInput{
		Limit: perPage,
		Page:  page,
	}

	if input.UserID.IsSet() {
		logsInput.ActorID = input.UserID.Value()
	}

	if input.ChannelID.IsSet() {
		logsInput.ChannelID = input.ChannelID.Value()
	}

	if input.ObjectID.IsSet() {
		logsInput.ObjectID = input.ObjectID.Value()
	}

	if input.System.IsSet() {
		systems := make([]string, 0, len(input.System.Value()))
		for _, s := range input.System.Value() {
			systems = append(systems, s.String())
		}

		logsInput.Systems = systems
	}

	if input.OperationType.IsSet() {
		operationTypes := make([]entity.AuditOperationType, 0, len(input.OperationType.Value()))
		for _, t := range input.OperationType.Value() {
			operationTypes = append(operationTypes, mappers.AuditTypeGqlToModel(t))
		}

		logsInput.OperationTypes = operationTypes
	}

	logs, err := r.deps.AuditLogsService.GetMany(ctx, logsInput)
	if err != nil {
		return nil, err
	}

	gqllogs := make([]gqlmodel.AdminAuditLog, 0, len(logs))
	for _, l := range logs {
		gqllogs = append(
			gqllogs,
			gqlmodel.AdminAuditLog{
				System:        mappers.AuditTableNameToGqlSystem(l.TableName),
				OperationType: mappers.AuditTypeModelToGql(l.OperationType),
				OldValue:      l.OldValue,
				NewValue:      l.NewValue,
				ObjectID:      l.ObjectID,
				UserID:        l.UserID,
				ChannelID:     l.ChannelID,
				CreatedAt:     l.CreatedAt,
			},
		)
	}

	total, err := r.deps.AuditLogsService.Count(ctx, audit_logs.GetCountInput{})
	if err != nil {
		return nil, err
	}

	return &gqlmodel.AdminAuditLogResponse{
		Logs:  gqllogs,
		Total: int(total),
	}, nil
}

// AdminAuditLog returns graph.AdminAuditLogResolver implementation.
func (r *Resolver) AdminAuditLog() graph.AdminAuditLogResolver { return &adminAuditLogResolver{r} }

type adminAuditLogResolver struct{ *Resolver }
