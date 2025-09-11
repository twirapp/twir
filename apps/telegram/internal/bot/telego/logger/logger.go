package logger

import (
	"fmt"

	"github.com/mymmrac/telego"
	twirlogger "github.com/twirapp/twir/libs/logger"
)

var _ telego.Logger = (*Logger)(nil)

func New(l twirlogger.Logger) *Logger {
	return &Logger{l}
}

type Logger struct {
	l twirlogger.Logger
}

func (l Logger) Debugf(format string, args ...any) {
	l.l.Debug(fmt.Sprintf(format, args))
}

func (l Logger) Errorf(format string, args ...any) {
	l.l.Error(fmt.Sprintf(format, args))
}
