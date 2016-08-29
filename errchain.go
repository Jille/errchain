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
func Chain(errs ...error) error {
	var ret []error
	for _, err := range errs {
		if err == nil {
			continue
		} else if es, ok := err.(errslice); ok {
			ret = append(ret, es.errors...)
		} else {
			ret = append(ret, err)
		}
	}
	switch len(ret) {
	case 0:
		return nil
	case 1:
		return ret[0]
	default:
		return errslice{ret}
	}
}

// Append changes err1 to be the combination of err1 and all others (nils allowed).
func Append(err1 *error, errs ...error) {
	*err1 = Chain(*err1, Chain(errs...))
}

// List turns an error in a list of errors.
func List(err error) []error {
	if err == nil {
		return nil
	}
	if e, ok := err.(errslice); ok {
		return e.errors
	}
	return []error{err}
}

// Call runs cb and chains the error to err. To be used from defers.
func Call(err *error, cb func() error) {
	*err = Chain(*err, cb())
}
