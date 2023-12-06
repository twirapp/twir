package helpers

import (
	"context"
	"fmt"
)

func GetSelectedDashboardIDFromCtx(ctx context.Context) (string, error) {
	dashboardId, ok := ctx.Value("dashboardId").(string)

	if !ok {
		return "", fmt.Errorf("failed to get dashboardId from context")
	}

	return dashboardId, nil
}
