package apperrors

import "errors"

func IsNotFound(err error) bool  { return errors.Is(err, ErrNotFound) }
func IsConflict(err error) bool  { return errors.Is(err, ErrConflict) }
func IsForbidden(err error) bool { return errors.Is(err, ErrForbidden) }
func IsRetryable(err error) bool {
	return errors.Is(err, ErrUnavailable) || errors.Is(err, ErrConflict)
}
func Code(err error) string {
	var c *Coded
	if errors.As(err, &c) {
		return c.Code
	}
	if errors.Is(err, ErrNotFound) {
		return "not_found"
	}
	if errors.Is(err, ErrConflict) {
		return "conflict"
	}
	if errors.Is(err, ErrForbidden) {
		return "forbidden"
	}
	return "internal"
}
