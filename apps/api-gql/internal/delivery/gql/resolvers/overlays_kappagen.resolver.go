package resolvers

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api-gql/internal/entity"
	"github.com/satont/twir/apps/api-gql/internal/gqlmodel"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/satont/twir/libs/logger/mappers"
)

func (r *queryResolver) OverlaysKappagen(ctx context.Context) (*gqlmodel.KappagenOverlay, error) {
	dashboardID, err := r.getDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	overlay, err := r.deps.KappagenService.GetOrCreate(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get kappagen overlay: %w", err)
	}

	return mapEntityToGQL(overlay), nil
}

func (r *mutationResolver) OverlaysKappagenUpdate(ctx context.Context, input gqlmodel.KappagenUpdateInput) (*gqlmodel.KappagenOverlay, error) {
	dashboardID, err := r.getDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := r.getUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get the current overlay for audit logging
	currentOverlay, err := r.deps.KappagenService.GetOrCreate(ctx, dashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current kappagen overlay: %w", err)
	}

	// Map input to entity
	entityInput := mapGQLInputToEntity(input, currentOverlay.ID)

	// Update the overlay
	updatedOverlay, err := r.deps.KappagenService.Update(ctx, dashboardID, entityInput)
	if err != nil {
		return nil, fmt.Errorf("failed to update kappagen overlay: %w", err)
	}

	// Log the audit
	r.deps.Logger.Audit(
		"Kappagen overlay update",
		audit.Fields{
			OldValue:      currentOverlay,
			NewValue:      updatedOverlay,
			ActorID:       lo.ToPtr(user.ID),
			ChannelID:     lo.ToPtr(dashboardID),
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelOverlayKappagen),
			OperationType: audit.OperationUpdate,
			ObjectID:      lo.ToPtr(updatedOverlay.ID.String()),
		},
	)

	return mapEntityToGQL(updatedOverlay), nil
}

// Helper functions to map between entity and GQL model
func mapEntityToGQL(entity entity.KappagenOverlay) *gqlmodel.KappagenOverlay {
	animations := make([]*gqlmodel.KappagenOverlayAnimationsSettings, 0, len(entity.Animations))
	for _, a := range entity.Animations {
		animations = append(animations, &gqlmodel.KappagenOverlayAnimationsSettings{
			Style: a.Style,
			Prefs: &gqlmodel.KappagenOverlayAnimationsPrefsSettings{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			},
			Count:   a.Count,
			Enabled: a.Enabled,
		})
	}

	return &gqlmodel.KappagenOverlay{
		ID:             entity.ID.String(),
		EnableSpawn:    entity.EnableSpawn,
		ExcludedEmotes: entity.ExcludedEmotes,
		EnableRave:     entity.EnableRave,
		Animation: &gqlmodel.KappagenOverlayAnimationSettings{
			FadeIn:  entity.Animation.FadeIn,
			FadeOut: entity.Animation.FadeOut,
			ZoomIn:  entity.Animation.ZoomIn,
			ZoomOut: entity.Animation.ZoomOut,
		},
		Animations: animations,
	}
}

func mapGQLInputToEntity(input gqlmodel.KappagenUpdateInput, id uuid.UUID) entity.KappagenOverlay {
	animations := make([]entity.KappagenOverlayAnimationsSettings, 0, len(input.Animations))
	for _, a := range input.Animations {
		animations = append(animations, entity.KappagenOverlayAnimationsSettings{
			Style: a.Style,
			Prefs: entity.KappagenOverlayAnimationsPrefsSettings{
				Size:    a.Prefs.Size,
				Center:  a.Prefs.Center,
				Speed:   a.Prefs.Speed,
				Faces:   a.Prefs.Faces,
				Message: a.Prefs.Message,
				Time:    a.Prefs.Time,
			},
			Count:   a.Count,
			Enabled: a.Enabled,
		})
	}

	return entity.KappagenOverlay{
		ID:             id,
		EnableSpawn:    input.EnableSpawn,
		ExcludedEmotes: input.ExcludedEmotes,
		EnableRave:     input.EnableRave,
		Animation: entity.KappagenOverlayAnimationSettings{
			FadeIn:  input.Animation.FadeIn,
			FadeOut: input.Animation.FadeOut,
			ZoomIn:  input.Animation.ZoomIn,
			ZoomOut: input.Animation.ZoomOut,
		},
		Animations: animations,
	}
}
