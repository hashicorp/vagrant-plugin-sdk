package errors

import (
	"fmt"
	"runtime"
	"strings"

	goerrors "errors"
)

var EnableCaller bool

type err struct {
	message  string
	function string
	path     string
	line     int
	wrap     error
}

func (e *err) Unwrap() error {
	return e.wrap
}

func (e *err) Error() string {
	msg := e.message
	if EnableCaller {
		msg = fmt.Sprintf("[%s@%s:%d] %s", e.function, e.path, e.line, e.message)
	}

	if e.wrap != nil {
		wMsg := e.wrap.Error()
		wErr, ok := e.wrap.(*err)
		if ok {
			wMsg = wErr.message
		}
		msg = fmt.Sprintf("%s %s", msg, wMsg)
	}

	return msg
}

func New(msg string) error {
	return new(msg)
}

func As(err error, target interface{}) bool {
	return goerrors.As(err, target)
}

func Is(err error, target error) bool {
	return goerrors.Is(err, target)
}

func Join(errs ...error) error {
	return goerrors.Join(errs...)
}

func Unwrap(err error) error {
	return goerrors.Unwrap(err)
}

func Wrap(msg string, originalErr error) error {
	newErr := new(msg).(*err)
	newErr.wrap = originalErr

	return newErr
}

// Create a new error with provided message. Will populate
// origin information if possible
func new(msg string) error {
	fn, file, line, ok := caller()
	if !ok {
		return &err{
			message: msg,
		}
	}

	return &err{
		message:  msg,
		function: fn,
		path:     file,
		line:     line,
	}
}

// Extracts caller information if possible and enabled
func caller() (fn string, file string, line int, ok bool) {
	if !EnableCaller {
		return
	}

	// NOTE: Actual caller will be three steps back
	pc, _, _, cok := runtime.Caller(3)
	if !cok {
		return
	}

	f := runtime.FuncForPC(pc)
	if f == nil {
		return
	}
	parts := strings.Split(f.Name(), "/")
	if len(parts) < 1 {
		return
	}
	fn = parts[len(parts)-1]
	file, line = f.FileLine(pc)
	ok = true

	return
}
