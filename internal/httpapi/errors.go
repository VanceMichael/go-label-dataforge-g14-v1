package httpapi

import (
	"encoding/json"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"net/http"
)

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Code: code, Message: message, RequestID: w.Header().Get("X-Request-ID")})
}
func StatusFor(err error) int {
	switch {
	case err == nil:
		return 200
	case apperrors.IsNotFound(err):
		return 404
	case apperrors.IsConflict(err):
		return 409
	case apperrors.IsForbidden(err):
		return 403
	default:
		return 500
	}
}
