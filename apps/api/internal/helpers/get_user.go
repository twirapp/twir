package helpers

import (
	"context"
	"fmt"

	model "github.com/twirapp/twir/libs/gomodels"
)

func GetUserModelFromCtx(ctx context.Context) (model.Users, error) {
	user, ok := ctx.Value("dbUser").(model.Users)
	if !ok {
		return model.Users{}, fmt.Errorf("failed to get user from context")
	}

	return user, nil
}
