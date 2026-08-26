package domain

import (
	"errors"
	"strings"
	"time"
)

var ErrEmptyName = errors.New("name is required")
var ErrInvalidCode = errors.New("code must be alphanumeric")
var ErrInvalidQuota = errors.New("quota must be positive")

func (r Resource) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.TenantID) == "" || strings.TrimSpace(r.OwnerID) == "" {
		return errors.New("resource identity is required")
	}
	if strings.TrimSpace(r.Name) == "" {
		return ErrEmptyName
	}
	if strings.TrimSpace(r.Code) == "" {
		return ErrInvalidCode
	}
	for _, c := range r.Code {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
			return ErrInvalidCode
		}
	}
	if len(r.Description) > 4000 {
		return errors.New("description too long")
	}
	return nil
}
func (a Authorization) Validate() error {
	if a.ID == "" || a.TenantID == "" || a.ResourceID == "" || a.ApplicantID == "" {
		return errors.New("authorization identity is required")
	}
	if strings.TrimSpace(a.Purpose) == "" {
		return errors.New("purpose is required")
	}
	if a.Quota <= 0 {
		return ErrInvalidQuota
	}
	return nil
}
func (l Lease) Validate() error {
	if l.ID == "" || l.AuthorizationID == "" || l.HolderID == "" {
		return errors.New("lease identity is required")
	}
	if l.ExpiresAt.IsZero() {
		return errors.New("lease expiry is required")
	}
	return nil
}
func (p Product) Validate() error {
	if p.ID == "" || p.TenantID == "" || p.OwnerID == "" || p.ResourceID == "" {
		return errors.New("product identity is required")
	}
	if strings.TrimSpace(p.Name) == "" {
		return ErrEmptyName
	}
	return nil
}
func (u User) Validate() error {
	if u.ID == "" || u.TenantID == "" {
		return errors.New("user identity is required")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email")
	}
	return nil
}
func (s Session) Valid(now time.Time) bool {
	return s.RevokedAt == nil && s.ExpiresAt != nil && s.ExpiresAt.After(now)
}
