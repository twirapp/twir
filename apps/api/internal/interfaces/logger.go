package interfaces

type Logger interface {
	Infow(msg string, keysAndValues ...any)
	Error(args ...any)
	Info(args ...any)
}
