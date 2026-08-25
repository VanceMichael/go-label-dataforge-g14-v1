package httpapi

import (
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/auth"
	"net/http"
)

type Authenticator struct{ Manager auth.Manager }

func (a Authenticator) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sid := r.Header.Get("X-Session-ID")
		token := r.Header.Get("X-Session-Token")
		if sid == "" || token == "" {
			writeErr(w, apperrors.ErrForbidden)
			return
		}
		u, e := a.Manager.Authenticate(r.Context(), sid, token)
		if e != nil {
			writeErr(w, e)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), u)))
	})
}
func RequireRole(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, ok := UserFrom(r.Context())
		if !ok || string(u.Role) != role {
			writeErr(w, apperrors.ErrForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
