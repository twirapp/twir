package users

import (
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func handleIgnoredUsersPost(services *types.Services, dto *ignoredUsersPostDto) error {
	for _, user := range dto.Users {
		newIgnoredUser := model.IgnoredUser{
			ID:          user.ID,
			Login:       null.StringFrom(user.Name),
			DisplayName: null.StringFrom(user.DisplayName),
		}
		err := services.Gorm.Save(&newIgnoredUser).Error
		if err != nil {
			services.Logger.Error(err)
			continue
		}
	}

	return nil
}
