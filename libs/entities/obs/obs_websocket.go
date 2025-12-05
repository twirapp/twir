package obs

type ObsWebsocketData struct {
	ID             int
	ChannelID      string
	ServerPort     int
	ServerAddress  string
	ServerPassword string
	Scenes         []string
	Sources        []string
	AudioSources   []string

	isNil bool
}

func (c ObsWebsocketData) IsNil() bool {
	return c.isNil
}

var NilObsWebsocket = ObsWebsocketData{
	isNil: true,
}

type ObsWebsocketCommandAction string

const (
	ObsWebsocketCommandActionSetScene       ObsWebsocketCommandAction = "SET_SCENE"
	ObsWebsocketCommandActionToggleSource   ObsWebsocketCommandAction = "TOGGLE_SOURCE"
	ObsWebsocketCommandActionToggleAudio    ObsWebsocketCommandAction = "TOGGLE_AUDIO"
	ObsWebsocketCommandActionSetVolume      ObsWebsocketCommandAction = "SET_VOLUME"
	ObsWebsocketCommandActionIncreaseVolume ObsWebsocketCommandAction = "INCREASE_VOLUME"
	ObsWebsocketCommandActionDecreaseVolume ObsWebsocketCommandAction = "DECREASE_VOLUME"
	ObsWebsocketCommandActionEnableAudio    ObsWebsocketCommandAction = "ENABLE_AUDIO"
	ObsWebsocketCommandActionDisableAudio   ObsWebsocketCommandAction = "DISABLE_AUDIO"
	ObsWebsocketCommandActionStartStream    ObsWebsocketCommandAction = "START_STREAM"
	ObsWebsocketCommandActionStopStream     ObsWebsocketCommandAction = "STOP_STREAM"
)

type ObsWebsocketCommand struct {
	Action      ObsWebsocketCommandAction
	Target      string
	VolumeValue *int
	VolumeStep  *int
}
