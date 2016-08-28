package errchain

import (
	"bytes"
	"fmt"
)

// errslice is a list of errors.
type errslice struct {
	errors []error
}

func (e errslice) Error() string {
	if len(e.errors) == 1 {
		return e.errors[0].Error()
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d errors: ", len(e.errors))
	for i, err := range e.errors {
		if i > 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// Chain takes two errors (or nils) and returns them combined if needed.
func Chain(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}
	var errs []error
	if err, ok := err1.(errslice); ok {
		errs = append(errs, err.errors...)
	} else {
		errs = append(errs, err1)
	}
	if err, ok := err2.(errslice); ok {
		errs = append(errs, err.errors...)
	} else {
		errs = append(errs, err2)
	}
	return errslice{errs}
}

// Append changes err1 to be the combination of err1 and err2 (nils allowed).
func Append(err1 *error, err2 error) {
	*err1 = Chain(*err1, err2)
}

// List turns an error in a list of errors.
func List(err error) []error {
	if e, ok := err.(errslice); ok {
		return e.errors
	}
	return []error{err}
}

// Call runs cb and chains the error to err. To be used from defers.
func Call(err *error, cb func() error) {
	*err = Chain(*err, cb())
}
