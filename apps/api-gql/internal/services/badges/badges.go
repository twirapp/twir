package badges

import (
	"context"

	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/badges"
	"github.com/twirapp/twir/libs/repositories/badges/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	BadgesRepository badges.Repository
	Config           config.Config
}

func New(opts Opts) *Service {
	return &Service{
		badgesRepository: opts.BadgesRepository,
		config:           opts.Config,
	}
}

type Service struct {
	badgesRepository badges.Repository
	config           config.Config
}

type GetManyInput struct {
	Enabled bool
}

func (c *Service) modelToEntity(b model.Badge) entity.Badge {
	return entity.Badge{
		ID:        b.ID,
		Name:      b.Name,
		Enabled:   b.Enabled,
		CreatedAt: b.CreatedAt,
		FileName:  b.FileName,
		FFZSlot:   b.FFZSlot,
		FileURL:   c.computeBadgeUrl(b.FileName),
	}
}

func (c *Service) computeBadgeUrl(fileName string) string {
	if c.config.AppEnv == "development" {
		return c.config.S3PublicUrl + "/" + c.config.S3Bucket + "/badges/" + fileName
	}

	return c.config.S3PublicUrl + "/badges/" + fileName
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.Badge, error) {
	selectedBadges, err := c.badgesRepository.GetMany(
		ctx,
		badges.GetManyInput{
			Enabled: input.Enabled,
		},
	)
	if err != nil {
		return nil, err
	}

	result := make([]entity.Badge, 0, len(selectedBadges))
	for _, b := range selectedBadges {
		result = append(result, c.modelToEntity(b))
	}

	return result, nil
}
