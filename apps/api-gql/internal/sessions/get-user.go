package sessions

import (
	"context"
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
)

func (s *Sessions) GetAuthenticatedUser(ctx context.Context) (*model.Users, error) {
	user, ok := s.sessionManager.Get(ctx, "dbUser").(model.Users)
	if !ok {
		return nil, fmt.Errorf("not authenticated")
	}

	return &user, nil
}
