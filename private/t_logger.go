package ekaweb_private

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Notice(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
	Crit(format string, args ...any)
	Alert(format string, args ...any)
	Emerg(format string, args ...any)
}
