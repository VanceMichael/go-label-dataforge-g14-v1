package httpapi

import (
	"encoding/json"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"net/http"
	"strings"
	"time"
)

func (h Handler) leases(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var l domain.Lease
	if json.NewDecoder(r.Body).Decode(&l) != nil {
		writeErr(w, errInvalid())
		return
	}
	if l.ID == "" {
		l.ID = time.Now().UTC().Format("20060102150405.000000000")
	}
	if l.Status == "" {
		l.Status = domain.LeaseReserved
	}
	l.CreatedAt = time.Now().UTC()
	if e := h.Sandbox.Reserve(r.Context(), l); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, l)
}
func (h Handler) leaseRun(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 4 {
		w.WriteHeader(404)
		return
	}
	var in struct{ Input string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, errInvalid())
		return
	}
	run, e := h.Sandbox.Run(r.Context(), parts[2], in.Input)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, run)
}
