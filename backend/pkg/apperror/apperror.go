// Package apperror defines typed errors that services return and handlers map to HTTP statuses.
package apperror

import "errors"

type Kind int

const (
	KindInternal Kind = iota
	KindValidation
	KindNotFound
	KindConflict
	KindUnauthorized
	KindForbidden
)

type Error struct {
	Kind Kind
	Msg  string
	Err  error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Msg + ": " + e.Err.Error()
	}
	return e.Msg
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Validation(msg string) *Error { return &Error{Kind: KindValidation, Msg: msg} }
func NotFound(msg string) *Error   { return &Error{Kind: KindNotFound, Msg: msg} }
func Conflict(msg string) *Error   { return &Error{Kind: KindConflict, Msg: msg} }
func Forbidden(msg string) *Error  { return &Error{Kind: KindForbidden, Msg: msg} }

func Unauthorized(msg string) *Error { return &Error{Kind: KindUnauthorized, Msg: msg} }

// Internal wraps an unexpected error; its cause is logged but never shown to clients.
func Internal(msg string, err error) *Error {
	return &Error{Kind: KindInternal, Msg: msg, Err: err}
}

// IsKind reports whether err (or anything it wraps) is an *Error of the given kind.
func IsKind(err error, kind Kind) bool {
	var appErr *Error
	return errors.As(err, &appErr) && appErr.Kind == kind
}
