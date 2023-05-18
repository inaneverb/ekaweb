package ekaweb_private

type Server interface {
	AsyncStart() error
	Stop() error
}
