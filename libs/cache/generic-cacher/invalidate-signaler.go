package generic_cacher

type InvalidateSignaler interface {
	Receiver() <-chan string
	Send(key string) error
}
