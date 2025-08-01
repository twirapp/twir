package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"time"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

// DashboardWidgetsEvents is the resolver for the dashboardWidgetsEvents field.
func (r *subscriptionResolver) DashboardWidgetsEvents(ctx context.Context) (<-chan *gqlmodel.DashboardEventListPayload, error) {
	dashboardID, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	channel := make(chan *gqlmodel.DashboardEventListPayload, 1)

	go func() {
		defer close(channel)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				events, err := r.deps.DashboardWidgetEventsService.GetDashboardWidgetsEvents(
					ctx,
					dashboardID,
					100,
				)
				if err != nil {
					r.deps.Logger.Error("cannot get dashboard events", err)
					time.Sleep(5 * time.Second)
					continue
				}

				mappedEvents := make([]gqlmodel.DashboardEventPayload, 0, len(events))

				for _, event := range events {
					mappedEvent, err := mappers.DashboardEventsDbToGql(event)
					if err != nil {
						r.deps.Logger.Error("cannot map dashboard event", err)
						continue
					}

					mappedEvents = append(mappedEvents, mappedEvent)
				}

				channel <- &gqlmodel.DashboardEventListPayload{
					Events: mappedEvents,
				}

				time.Sleep(5 * time.Second)
			}
		}
	}()

	return channel, nil
}
