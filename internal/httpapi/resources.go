package httpapi

import (
	"encoding/json"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"net/http"
	"strings"
	"time"
)

type resourceInput struct{ TenantID, OwnerID, Code, Name, Description string }

func (h Handler) resources(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var in resourceInput
		if json.NewDecoder(r.Body).Decode(&in) != nil {
			writeErr(w, errInvalid())
			return
		}
		now := time.Now().UTC()
		res := domain.Resource{ID: now.Format("20060102150405.000000000"), TenantID: in.TenantID, OwnerID: in.OwnerID, Code: in.Code, Name: in.Name, Description: in.Description, Status: domain.ResourceDraft, CreatedAt: now, UpdatedAt: now}
		ver := domain.ResourceVersion{ID: res.ID + "-v1", ResourceID: res.ID, Number: 1, Schema: "json", ContentHash: "pending", CreatedAt: now}
		rev := domain.Review{ID: res.ID + "-review", ResourceID: res.ID, Status: "pending", Notes: "", CreatedAt: now}
		created, e := h.Registry.Create(r.Context(), res, ver, rev, r.Header.Get("X-Request-ID"))
		if e != nil {
			writeErr(w, e)
			return
		}
		writeJSON(w, 201, created)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}
func (h Handler) resourceAction(w http.ResponseWriter, r *http.Request) {
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
	var in struct {
		Version           int
		TenantID, ActorID string
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, errInvalid())
		return
	}
	if e := h.Registry.Submit(r.Context(), id, in.Version, in.ActorID, in.TenantID, r.Header.Get("X-Request-ID")); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "submitted", "id": id})
}
