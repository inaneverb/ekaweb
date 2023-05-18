package ekaweb_respondent

// Expander is an interface that is a part of Respondent middleware.
// It allows you to get a Manifest from the error object
// that will provide detailed description of the occurred error.
//
// Technically an implementation of Expander must recognize what kind of error
// is passed, parse it, match it with the pattern and expand an error.
// You may use CommonExpander as a default implementation of Expander interface.
type Expander interface {
	Expand(err error) *Manifest
}

// Replacer is an interface that is a part of Respondent middleware.
// It allows you to switch one error to another.
type Replacer interface {
	Replace(err error) error
}

// Applicator is an interface that is a part of Respondent middleware.
// It allows you to transform a Manifest obtained by an Expander to an HTTP response.
type Applicator interface {
	Apply(ctx any, manifest *Manifest)
}

// Manifest is a structured representation of error, the main goal of which
// is to represent an occurred error as an HTTP response.
type Manifest struct {
	Status    int
	Error     string
	ErrorID   string
	ErrorCode int

	ErrorDetail  string
	ErrorDetails []string

	customFillers []ManifestCustomFiller
}

type ManifestCustomFiller = func(ctx any, manifest *Manifest)

// Clone returns a full copy of the Manifest.
func (m *Manifest) Clone() *Manifest {

	if m != nil {
		m2 := *m
		m = &m2
	}

	return m
}
