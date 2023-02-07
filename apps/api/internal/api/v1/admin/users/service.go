package users

import (
	"github.com/guregu/null"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
	"gorm.io/gorm"
)

func handleIgnoredUsersPost(services types.Services, dto *ignoredUsersPostDto) error {
	db := do.MustInvoke[*gorm.DB](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	for _, user := range dto.Users {
		newIgnoredUser := model.IgnoredUser{
			ID:          user.ID,
			Login:       null.StringFrom(user.Name),
			DisplayName: null.StringFrom(user.DisplayName),
		}
		err := db.Save(&newIgnoredUser).Error
		if err != nil {
			logger.Error(err)
			continue
		}
	}

	return nil
}
