package httpapi

import (
	"encoding/json"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"net/http"
	"strings"
	"time"
)

func (h Handler) products(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var p domain.Product
	if json.NewDecoder(r.Body).Decode(&p) != nil {
		writeErr(w, errInvalid())
		return
	}
	if p.ID == "" {
		p.ID = time.Now().UTC().Format("20060102150405.000000000")
	}
	if p.Status == "" {
		p.Status = domain.ProductDraft
	}
	p.CreatedAt = time.Now().UTC()
	p.UpdatedAt = p.CreatedAt
	if e := h.Products.Create(r.Context(), p); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, p)
}
func (h Handler) productAction(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 3 {
		w.WriteHeader(404)
		return
	}
	var in struct{ Version int }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, errInvalid())
		return
	}
	if e := h.Products.Publish(r.Context(), parts[2], in.Version); e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 200, map[string]string{"status": "published", "id": parts[2]})
}
