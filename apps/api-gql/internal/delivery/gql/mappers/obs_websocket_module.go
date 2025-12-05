package mappers
import (
"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
obsentity "github.com/twirapp/twir/libs/entities/obs"
)
func MapObsWebsocketModuleDataToGql(
data *obsentity.ObsWebsocketData,
) *gqlmodel.ObsWebsocketModule {
	if data == nil {
		return nil
	}
	return &gqlmodel.ObsWebsocketModule{
		ServerPort:     data.ServerPort,
		ServerAddress:  data.ServerAddress,
		ServerPassword: data.ServerPassword,
		Sources:        data.Sources,
		AudioSources:   data.AudioSources,
		Scenes:         data.Scenes,
	}
}
func MapObsWebsocketCommandToGql(
cmd *obsentity.ObsWebsocketCommand,
) *gqlmodel.ObsWebsocketCommand {
	if cmd == nil {
		return nil
	}
	return &gqlmodel.ObsWebsocketCommand{
		Action:      gqlmodel.ObsWebsocketCommandAction(cmd.Action),
		Target:      cmd.Target,
		VolumeValue: cmd.VolumeValue,
		VolumeStep:  cmd.VolumeStep,
	}
}
func MapGqlActionToEntity(action gqlmodel.ObsWebsocketCommandAction) obsentity.ObsWebsocketCommandAction {
	return obsentity.ObsWebsocketCommandAction(action)
}
