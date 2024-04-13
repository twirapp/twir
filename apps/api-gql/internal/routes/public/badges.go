package public

import (
	"github.com/gin-gonic/gin"
	model "github.com/satont/twir/libs/gomodels"
)

func (p *Public) computeBadgeUrl(fileName string) string {
	if p.config.AppEnv == "development" {
		return p.config.S3PublicUrl + "/" + p.config.S3Bucket + "/badges/" + fileName
	}

	return p.config.S3Host + "/badges/" + fileName
}

func (p *Public) HandleBadgesGet(c *gin.Context) {
	var badges []model.Badge
	if err := p.gorm.
		WithContext(c.Request.Context()).
		Find(&badges).
		Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	mappedBadges := make([]map[string]any, 0, len(badges))

	for _, badge := range badges {
		mappedBadges = append(
			mappedBadges,
			map[string]any{
				"name":    badge.Name,
				"url":     p.computeBadgeUrl(badge.FileName),
				"ffzSlot": badge.FFZSlot,
				"users":   badge.Users,
			},
		)
	}

	c.JSON(200, mappedBadges)
}
