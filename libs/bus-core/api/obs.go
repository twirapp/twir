package api

const (
	TriggerObsCommandSubject = "api.obs.command.trigger"
)

type ObsCommandAction string

const (
	ObsCommandActionSetScene       ObsCommandAction = "setScene"
	ObsCommandActionToggleSource   ObsCommandAction = "toggleSource"
	ObsCommandActionToggleAudio    ObsCommandAction = "toggleAudioSource"
	ObsCommandActionSetVolume      ObsCommandAction = "setVolume"
	ObsCommandActionIncreaseVolume ObsCommandAction = "increaseVolume"
	ObsCommandActionDecreaseVolume ObsCommandAction = "decreaseVolume"
	ObsCommandActionEnableAudio    ObsCommandAction = "enableAudio"
	ObsCommandActionDisableAudio   ObsCommandAction = "disableAudio"
	ObsCommandActionStartStream    ObsCommandAction = "startStream"
	ObsCommandActionStopStream     ObsCommandAction = "stopStream"
)

type TriggerObsCommand struct {
	ChannelId   string
	Action      ObsCommandAction
	Target      string // scene name, source name, or audio source name
	VolumeValue *int   // optional, only for setVolume
	VolumeStep  *int   // optional, only for increaseVolume/decreaseVolume
}
