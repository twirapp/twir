package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

var commandsExpiresAtMap = map[model.ChannelCommandExpiresType]gqlmodel.CommandExpiresType{
	model.ChannelCommandExpiresTypeDelete:  gqlmodel.CommandExpiresTypeDelete,
	model.ChannelCommandExpiresTypeDisable: gqlmodel.CommandExpiresTypeDisable,
}

func CommandsExpiresAtDbToGql(in model.ChannelCommandExpiresType) gqlmodel.CommandExpiresType {
	return commandsExpiresAtMap[in]
}

func CommandsExpiresAtGqlToDb(in gqlmodel.CommandExpiresType) model.ChannelCommandExpiresType {
	for k, v := range commandsExpiresAtMap {
		if v == in {
			return k
		}
	}

	return model.ChannelCommandExpiresTypeDelete
}
