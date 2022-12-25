package interfaces

type Logger interface {
	Infow(msg string, args ...any)
	Error(args ...any)
	Info(args ...any)
}
