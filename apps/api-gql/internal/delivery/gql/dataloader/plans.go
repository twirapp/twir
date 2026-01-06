package dataloader

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func (c *dataLoader) getPlansByChannelIDs(ctx context.Context, inputs []GetPlayByChannelId) (
	[]*gqlmodel.Plan,
	[]error,
) {
	uniquePlanIDsMap := make(map[string]struct{})
	for _, i := range inputs {
		if i.PlanID == nil || *i.PlanID == "" {
			continue
		}

		if _, ok := uniquePlanIDsMap[*i.PlanID]; ok {
			continue
		}

		uniquePlanIDsMap[*i.PlanID] = struct{}{}
	}

	uniquePlanIDsSlice := make([]string, 0, len(uniquePlanIDsMap))
	for planID := range uniquePlanIDsMap {
		uniquePlanIDsSlice = append(uniquePlanIDsSlice, planID)
	}

	plans, err := c.deps.PlansRepository.GetManyByIDs(ctx, uniquePlanIDsSlice)
	if err != nil {
		return nil, []error{err}
	}

	plansByID := make(map[string]*gqlmodel.Plan)
	for _, p := range plans {
		if !p.IsNil() {
			plansByID[p.ID] = mappers.PlanToGql(p)
		}
	}

	result := make([]*gqlmodel.Plan, len(inputs))
	for i, ch := range inputs {
		if ch.PlanID == nil || *ch.PlanID == "" {
			continue
		}

		plan, ok := plansByID[*ch.PlanID]
		if !ok {
			continue
		}

		result[i] = plan
	}

	return result, nil
}

type GetPlayByChannelId struct {
	ChannelID string
	PlanID    *string
}

func GetPlanByChannelID(ctx context.Context, input GetPlayByChannelId) (*gqlmodel.Plan, error) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.plansByChannelIDLoader.Load(ctx, input)
}

func GetPlansByChannelIDs(ctx context.Context, inputs []GetPlayByChannelId) (
	[]*gqlmodel.Plan,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.plansByChannelIDLoader.LoadAll(ctx, inputs)
}
