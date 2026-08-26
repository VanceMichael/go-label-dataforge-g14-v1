package apperrors

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrConflict    = errors.New("conflict")
	ErrForbidden   = errors.New("forbidden")
	ErrInvalid     = errors.New("invalid request")
	ErrExpired     = errors.New("expired")
	ErrUnavailable = errors.New("dependency unavailable")
)

type Coded struct {
	Code string
	Err  error
}

func (e *Coded) Error() string { return e.Code + ": " + e.Err.Error() }
func (e *Coded) Unwrap() error { return e.Err }
func Wrap(code string, err error) error {
	if err == nil {
		return nil
	}
	return &Coded{Code: code, Err: err}
}
