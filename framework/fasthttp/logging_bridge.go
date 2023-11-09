package ekaweb_fasthttp

import (
	"github.com/inaneverb/ekaweb/v2"
	"github.com/valyala/fasthttp"
)

type zapLoggerAsFasthttpLogger struct {
	origin ekaweb.Logger
}

func (z zapLoggerAsFasthttpLogger) Printf(format string, args ...any) {
	z.origin.Error(format, args...)
}

func newLoggingBridge(log ekaweb.Logger) fasthttp.Logger {
	return zapLoggerAsFasthttpLogger{log}
}
