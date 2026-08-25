package httpapi

import (
	"encoding/json"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"net/http"
	"strings"
	"time"
)

func (h Handler) authorizations(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var a domain.Authorization
	if json.NewDecoder(r.Body).Decode(&a) != nil {
		writeErr(w, errInvalid())
		return
	}
	if a.ID == "" {
		a.ID = time.Now().UTC().Format("20060102150405.000000000")
	}
	a.Status = domain.AuthRequested
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = a.CreatedAt
	if e := h.Authorizer.Request(r.Context(), a); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, a)
}
func (h Handler) authorizationAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		w.WriteHeader(404)
		return
	}
	id := parts[2]
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct{ Version int }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, errInvalid())
		return
	}
	if e := h.Authorizer.Approve(r.Context(), id, in.Version); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "active", "id": id})
}
