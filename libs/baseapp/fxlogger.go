package baseapp

import (
	"reflect"

	"go.uber.org/fx/fxevent"
)

type errorOnlyFxLogger struct {
	fxevent.ConsoleLogger
}

// LogEvent filters to only log if the event has a non-nil Err field.
func (l *errorOnlyFxLogger) LogEvent(event fxevent.Event) {
	v := reflect.ValueOf(event)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	errField := v.FieldByName("Err")
	if errField.IsValid() && !errField.IsNil() {
		if err, ok := errField.Interface().(error); ok && err != nil {
			l.ConsoleLogger.LogEvent(event)
		}
	}
}
