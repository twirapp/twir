package interfaces

type Logger interface {
	Error(args ...any)
	Infow(msg string, keysAndValues ...interface{})
	Info(args ...interface{})
}
