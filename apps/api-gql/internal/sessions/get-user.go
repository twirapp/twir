package sessions

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
)

func GetAuthenticatedUser(ctx context.Context) (*model.Users, error) {
	fmt.Println("2")
	user, ok := ctx.Value("dbUser").(model.Users)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	return &user, nil
}
