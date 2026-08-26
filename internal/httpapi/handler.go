package httpapi

import (
	"encoding/json"
	"errors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/auth"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/service"
	"net/http"
	"strings"
)

type Handler struct {
	Auth       auth.Manager
	Registry   service.Registry
	Authorizer service.Authorizer
	Sandbox    service.Sandbox
	Products   service.ProductService
}

func New(h Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	})
	mux.HandleFunc("/v1/sessions", h.sessions)
	mux.HandleFunc("/v1/sessions/", h.logout)
	mux.HandleFunc("/v1/resources", h.resources)
	mux.HandleFunc("/v1/resources/", h.resourceAction)
	mux.HandleFunc("/v1/authorizations", h.authorizations)
	mux.HandleFunc("/v1/authorizations/", h.authorizationAction)
	mux.HandleFunc("/v1/leases", h.leases)
	mux.HandleFunc("/v1/leases/", h.leaseRun)
	mux.HandleFunc("/v1/products", h.products)
	mux.HandleFunc("/v1/products/", h.productAction)
	return Recover(RequestID(mux))
}
func (h Handler) sessions(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		return
	}
	var in struct{ Tenant, Email string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, apperrors.ErrInvalid)
		return
	}
	s, u, e := h.Auth.Login(r.Context(), in.Tenant, in.Email)
	if e != nil {
		writeErr(w, e)
		return
	}
	writeJSON(w, 201, map[string]any{"session": s, "user": u})
}
func (h Handler) logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(405)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/v1/sessions/")
	if e := h.Auth.Logout(r.Context(), id); e != nil {
		writeErr(w, e)
		return
	}
	w.WriteHeader(204)
}
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, e error) {
	status := 500
	code := "internal"
	switch {
	case errors.Is(e, apperrors.ErrNotFound):
		status = 404
		code = "not_found"
	case errors.Is(e, apperrors.ErrConflict):
		status = 409
		code = "conflict"
	case errors.Is(e, apperrors.ErrForbidden):
		status = 403
		code = "forbidden"
	case errors.Is(e, apperrors.ErrInvalid):
		status = 400
		code = "invalid"
	case errors.Is(e, apperrors.ErrExpired):
		status = 410
		code = "expired"
	}
	writeJSON(w, status, map[string]string{"code": code, "message": e.Error()})
}
func errInvalid() error { return apperrors.ErrInvalid }
