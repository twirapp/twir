package donatepay

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/api/v1/integrations/helpers"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"github.com/satont/tsuwari/apps/api/internal/types"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func handleGet(services types.Services, channelId string) (*model.ChannelsIntegrationsData, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	integration, err := helpers.GetIntegration(channelId, "DONATEPAY", services.DB)
	if err != nil {
		logger.Error(err)
		return nil, nil
	}

	if integration == nil {
		return nil, nil
	}

	return integration.Data, nil
}
