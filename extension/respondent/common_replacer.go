package ekaweb_respondent

import (
	"errors"

	"github.com/inaneverb/ekacore/ekaunsafe/v4"
)

// CommonReplacer is an error "replacer". Implements Replacer interface.
// The core of the logic is Replace() method. It allows you to replace one error
// by the another.
//
// You must use ReplaceBy() method to make a rule what error must be replaced
// by which. The Replace() method is not for you!
//
// It's ready-to-use object after instantiating this type but feel free to use
// its constructor: NewCommonReplacer().
type CommonReplacer struct {
	direct map[error]error
	deep   []commonReplacerDeepPair
	check  []commonReplacerCheckPair
	custom []func(err error) error
}

// commonReplacerDeepPair is a unit that contains:
//  1. An original error that will be used at the error pattern matching
//  2. A replacement error that will be used instead of original
//     if pattern matching is succeeded.
type commonReplacerDeepPair struct {
	Original    error
	Replacement error
}

// commonReplacerCheckPair is a unit that contains:
//  1. A callback that will tell whether provided error is that one,
//     that should be replaced by the replacement.
//  2. A replacement error that will be used instead of original.
type commonReplacerCheckPair struct {
	Checker     func(err error) bool
	Replacement error
}

// Replace returns an error that is associated with the provided
// and provided error is exchanged to. Implements Replacer interface.
//
// If replacement rule is not registered for the provided error
// it remains intact and is returned as is.
func (cr *CommonReplacer) Replace(err error) error {

	// "if err == nil" statement reports false when err is (*someUserErrorType)(nil)
	// for example. So, we're using more explicit check here.

	if cr == nil || ekaunsafe.UnpackInterface(err).Word == nil {
		return nil
	}

	var iErr = ekaunsafe.UnpackInterface(err)
	var errKind = ekaunsafe.ReflectTypeOfRType(iErr.Type).Kind()

	if isGoHashableObject(errKind) {
		if replacement, found := cr.direct[err]; found {
			return replacement
		}
	}

	for i, n := 0, len(cr.deep); i < n; i++ {
		if errors.Is(err, cr.deep[i].Original) {
			return cr.deep[i].Replacement
		}
	}

	for i, n := 0, len(cr.check); i < n; i++ {
		if cr.check[i].Checker(err) {
			return cr.check[i].Replacement
		}
	}

	for i, n := 0, len(cr.custom); i < n; i++ {
		if err := cr.custom[i](err); err != nil {
			return err
		}
	}

	return err
}

// ReplaceBy creates a new replacement rule: the 'original' error should be replaced
// by the 'replacement' error, maybe with check if such param is specified.
//
// Enabling deep scanning means that error's pattern matching at Replace()
// will be performed by go1.13 API errors.Is() instead of direct comparison.
//
// Does nothing if any of error is nil, or they are the same.
// It's ok if 'replacement' is nil. It will lead to getting ErrEmptyReplacement
// by the Respondent.
func (cr *CommonReplacer) ReplaceBy(original, replacement error, deep ...bool) *CommonReplacer {

	deepValue := len(deep) > 0 && deep[0]

	if cr == nil ||
		ekaunsafe.UnpackInterface(original).Word == nil ||
		(!deepValue && original == replacement) ||
		(deepValue && (errors.Is(original, replacement) || errors.Is(replacement, original))) {

		return cr
	}

	if deepValue {
		cr.deep = append(cr.deep, commonReplacerDeepPair{original, replacement})
	} else {
		if cr.direct == nil {
			cr.direct = make(map[error]error)
		}
		cr.direct[original] = replacement
	}

	return cr
}

// ReplaceByChecker creates a new replacement rule: the original error should be replaced
// if provided callback (the original error is passed to) will return true.
func (cr *CommonReplacer) ReplaceByChecker(checker func(err error) bool, replacement error) *CommonReplacer {
	if cr != nil && checker != nil {
		cr.check = append(cr.check, commonReplacerCheckPair{checker, replacement})
	}
	return cr
}

// ReplaceByCustom creates a new replacement rule: the original error should be replaced
// by that one that is returned by the provided callback. The original error will be
// passed to that callback. If callback returns nil the flow continues processing.
// If you need to replace error by nil, use ReplaceByChecker.
func (cr *CommonReplacer) ReplaceByCustom(custom func(original error) error) *CommonReplacer {
	if cr != nil && custom != nil {
		cr.custom = append(cr.custom, custom)
	}
	return cr
}

// ReplaceByBulk allows you to create and specify replacement rules by more convenient way
// using map instead of chaining ReplaceBy() calls.
// The provided map will be consumed and all errors will be copied to call
// ReplaceBy() under the hood. So this method is just helper.
//
// In the passed map the key must be an original error and the value should be
// the replacement.
func (cr *CommonReplacer) ReplaceByBulk(rules map[error]error, deep ...bool) *CommonReplacer {

	// No need for any nil check. ReplaceBy has the all necessary checks,
	// and it's OK to iterate over nil map.

	for original, replacement := range rules {
		cr.ReplaceBy(original, replacement, deep...)
	}

	return cr
}

// NewCommonReplacer is a constructor of CommonReplacer. It initializes all
// internal fields and returns a ready-to-use object.
func NewCommonReplacer() *CommonReplacer {
	return &CommonReplacer{
		direct: make(map[error]error),
		deep:   make([]commonReplacerDeepPair, 0, 32),
		check:  make([]commonReplacerCheckPair, 0, 32),
	}
}
