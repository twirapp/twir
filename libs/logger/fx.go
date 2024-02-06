package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

var FxDiOnlyErrors = fx.WithLogger(
	func() fxevent.Logger {
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
		z := zap.Must(config.Build())

		return &fxevent.ZapLogger{Logger: z}
	},
)
