package health

import (
	"context"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/storage/sqlite"
	"net/http"
)

type Handler struct{ DB *sqlite.DB }

func (h Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
func (h Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if e := h.DB.Ping(r.Context()); e != nil {
		http.Error(w, `{"status":"not_ready"}`, http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ready"}`))
}

var _ = context.Background
