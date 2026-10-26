package domain

import "errors"

// Ledger rule violations. Repositories return these; services map them to
// apperror kinds so the HTTP layer can pick the right status.
var (
	ErrFundNotFound          = errors.New("fund not found")
	ErrInsufficientBalance   = errors.New("insufficient fund balance")
	ErrAlreadyReversed       = errors.New("transaction has already been reversed")
	ErrCannotReverseReversal = errors.New("a reversal entry cannot itself be reversed")
)

// Qurban booking rule violations.
var (
	ErrPackageNotFound   = errors.New("qurban package not found")
	ErrPackageOutOfStock = errors.New("qurban package is out of stock")
)
