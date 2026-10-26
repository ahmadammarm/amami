// Package apperror defines typed errors that services return and handlers map to HTTP statuses.
package apperror

import (
	"errors"

	"gorm.io/gorm"
)

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

// FromStore converts an error returned by a repository into an *Error:
//   - an *Error (typically a business rule rejected inside a transaction)
//     passes through unchanged;
//   - a missing row becomes NotFound(notFoundMsg);
//   - anything else becomes Internal(internalMsg), keeping the cause for logs.
//
// Callers check for other specific database errors (duplicate key, foreign
// key) before calling it when those need their own message.
func FromStore(err error, notFoundMsg, internalMsg string) error {
	if err == nil {
		return nil
	}
	var appErr *Error
	if errors.As(err, &appErr) {
		return appErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound(notFoundMsg)
	}
	return Internal(internalMsg, err)
}

// IsKind reports whether err (or anything it wraps) is an *Error of the given kind.
func IsKind(err error, kind Kind) bool {
	var appErr *Error
	return errors.As(err, &appErr) && appErr.Kind == kind
}
