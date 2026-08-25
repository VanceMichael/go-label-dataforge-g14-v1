package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
	"strings"
	"time"
)

type Manager struct {
	Store repository.Store
	Clock clock.Clock
	TTL   time.Duration
}

func (m Manager) Login(ctx context.Context, tenant, email string) (domain.Session, domain.User, error) {
	u, e := m.Store.FindUserByEmail(ctx, tenant, email)
	if e != nil {
		return domain.Session{}, u, e
	}
	if !u.Active {
		return domain.Session{}, u, apperrors.ErrForbidden
	}
	raw := make([]byte, 32)
	if _, e = rand.Read(raw); e != nil {
		return domain.Session{}, u, e
	}
	now := m.Clock.Now()
	exp := now.Add(m.TTL)
	id := hex.EncodeToString(raw[:8])
	token := hex.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	s := domain.Session{ID: id, UserID: u.ID, TokenHash: hex.EncodeToString(sum[:]), ExpiresAt: &exp, CreatedAt: &now}
	return s, u, m.Store.CreateSession(ctx, s)
}
func (m Manager) Authenticate(ctx context.Context, id, token string) (domain.User, error) {
	s, e := m.Store.FindSession(ctx, id)
	if e != nil {
		return domain.User{}, e
	}
	if s.RevokedAt != nil || s.ExpiresAt == nil || !s.ExpiresAt.After(m.Clock.Now()) {
		return domain.User{}, apperrors.ErrForbidden
	}
	sum := sha256.Sum256([]byte(token))
	if !strings.EqualFold(hex.EncodeToString(sum[:]), s.TokenHash) {
		return domain.User{}, apperrors.ErrForbidden
	}
	return domain.User{ID: s.UserID}, nil
}
func (m Manager) Logout(ctx context.Context, id string) error { return m.Store.RevokeSession(ctx, id) }
