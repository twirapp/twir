package timers

const (
	AddTimerSubject    = "timers.add"
	RemoveTimerSubject = "timers.remove"
)

type AddOrRemoveTimerRequest struct {
	TimerID string
}
