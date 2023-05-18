package ekaweb_nbio

import (
	"github.com/inaneverb/ekaweb"
	"github.com/lesismal/nbio/logging"
)

type zapLoggerAsNbioLogger struct {
	origin ekaweb.Logger
}

func (z zapLoggerAsNbioLogger) SetLevel(_ int) {
	// Do nothing. Nbio package does not call this method.
}

func (z zapLoggerAsNbioLogger) Debug(format string, v ...any) {
	z.origin.Debug(format, v...)
}

func (z zapLoggerAsNbioLogger) Info(format string, v ...any) {
	z.origin.Info(format, v...)
}

func (z zapLoggerAsNbioLogger) Warn(format string, v ...any) {
	z.origin.Warn(format, v...)
}

func (z zapLoggerAsNbioLogger) Error(format string, v ...any) {
	z.origin.Error(format, v...)
}

func newLoggingBridge(log ekaweb.Logger) logging.Logger {
	return zapLoggerAsNbioLogger{log}
}
