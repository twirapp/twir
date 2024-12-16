package public

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	badges_with_users "github.com/twirapp/twir/apps/api-gql/internal/services/badges-with-users"
)

func (p *Public) computeBadgeUrl(fileName string) string {
	if p.config.AppEnv == "development" {
		return p.config.S3PublicUrl + "/" + p.config.S3Bucket + "/badges/" + fileName
	}

	return p.config.S3PublicUrl + "/badges/" + fileName
}

type badgeWithUsers struct {
	ID      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	FFZSlot int       `json:"ffzSlot"`
	URL     string    `json:"url"`
	Users   []string  `json:"users"`
}

func (p *Public) HandleBadgesGet(c *gin.Context) {
	entities, err := p.badgesWithUsersService.GetMany(
		c.Request.Context(),
		badges_with_users.GetManyInput{Enabled: true},
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result := make([]badgeWithUsers, 0, len(entities))
	for _, entity := range entities {
		result = append(
			result,
			badgeWithUsers{
				ID:      entity.ID,
				Name:    entity.Name,
				FFZSlot: entity.FFZSlot,
				URL:     entity.FileURL,
				Users:   entity.Users,
			},
		)
	}

	c.JSON(200, result)
}
