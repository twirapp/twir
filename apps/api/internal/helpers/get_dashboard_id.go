package helpers

import (
	"context"
	"fmt"
)

var ErrDashboardIdNotFound = fmt.Errorf("failed to get dashboardId from context")

func GetSelectedDashboardIDFromContext(ctx context.Context) (string, error) {
	dashboardId, ok := ctx.Value("dashboardId").(string)
	if !ok {
		return "", ErrDashboardIdNotFound
	}

	if dashboardId == "" {
		return "", fmt.Errorf("dashboardId is empty")
	}

	return dashboardId, nil
}
