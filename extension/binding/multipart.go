package ekaweb_bind

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/inaneverb/ekaweb/v2"
)

//goland:noinspection GoErrorStringFormat
var (
	ErrMultipartIncorrectMIME = errors.New("Multipart: Incorrect MIME")
	ErrMultipartIncorrectFile = errors.New("Multipart: Incorrect file")
	ErrMultipartFileNotFound  = errors.New("Multipart: File not found")
)

// FormDataFile is like sequenced
// http.Request.ParseMultipartForm() + http.Request.FormFile(),
// but w/o their limitations and more convenient signature.
//
// Rules:
//
//   - You don't have to specify filename (`name` parameter).
//     Don't know which filename user will use? No problem, get exactly that one
//     that used provides (if only one was provided, error otherwise).
//
//   - You don't have to specify max available RAM to parse (`maxMemory` parameter).
//     Don't know which value is better? No problem, Golang default is used.
//
//   - Only one call. You even don't have to check errors. All of them already saved,
//     using ErrorApply() and ErrorDetailApply(). Just check returned arguments.
//     If they are nil, an error occurred, and it's already saved to user context.
//
// Returned errors:
// - ErrMultipartIncorrectMIME: Incorrect MIME type or data is empty;
// - ErrMultipartIncorrectFile: File too large, file is empty, RAM errors, etc;
// - ErrMultipartFileNotFound: Filename is specified, but not found.
func FormDataFile(
	r *http.Request, name string,
	maxMemory int64) (multipart.File, *multipart.FileHeader) {

	multipartReader, err := r.MultipartReader()
	switch {
	case err != nil:
		ekaweb.ErrorDetailApply(r, err.Error())
		fallthrough

	case multipartReader == nil:
		ekaweb.ErrorApply(r, ErrMultipartIncorrectMIME)
		return nil, nil
	}

	if maxMemory <= 0 {
		maxMemory = 32 << 20 // http.defaultMaxMemory (32 MB)
	}

	form, err := multipartReader.ReadForm(maxMemory)
	switch {
	case err != nil:
		ekaweb.ErrorDetailApply(r, err.Error())
		fallthrough

	case form == nil || len(form.File) == 0:
		ekaweb.ErrorApply(r, ErrMultipartIncorrectFile)
		return nil, nil
	}

	// Ok, the last thing. It's only file we're interested in.
	// If name is not specified, it's our responsibility to figure it out.

	switch {
	case name == "" && len(form.File) > 1:
		ekaweb.ErrorApply(r, ErrMultipartFileNotFound)
		return nil, nil

	case name == "":
		for fileName := range form.File {
			name = fileName
		}
	}

	fileHeaders := form.File[name]
	if len(fileHeaders) == 0 {
		ekaweb.ErrorApply(r, ErrMultipartFileNotFound)
		return nil, nil
	}

	f, err := fileHeaders[0].Open()
	if err != nil {
		ekaweb.ErrorApply(r, ErrMultipartIncorrectFile)
		ekaweb.ErrorDetailApply(r, err.Error())
		return nil, nil
	}

	return f, fileHeaders[0]
}
