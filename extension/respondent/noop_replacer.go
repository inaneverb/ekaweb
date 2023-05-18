package ekaweb_respondent

// noOpReplacer is an error "replacer". Implements Replacer interface.
// The core method Replace() does absolutely nothing but returns a passed error.
type noOpReplacer struct{}

// Replace returns a provided error without any modifications.
func (*noOpReplacer) Replace(err error) error {
	return err
}

func newNoOpReplacer() *noOpReplacer {
	return new(noOpReplacer)
}
