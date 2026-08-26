package httpapi

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"net/http"
)

type key string

const userKey key = "user"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			id = hex.EncodeToString(b)
		}
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r)
	})
}
func WithUser(ctx context.Context, u domain.User) context.Context {
	return context.WithValue(ctx, userKey, u)
}
func UserFrom(ctx context.Context) (domain.User, bool) {
	u, ok := ctx.Value(userKey).(domain.User)
	return u, ok
}
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recover() != nil {
				http.Error(w, `{"error":"internal"}`, 500)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
