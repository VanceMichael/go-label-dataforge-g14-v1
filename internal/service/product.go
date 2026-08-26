package service

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/apperrors"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/clock"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/repository"
)

type ProductService struct {
	Store repository.Store
	Clock clock.Clock
}

func (s ProductService) Create(ctx context.Context, p domain.Product) error {
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.CreateProductTx(ctx, tx, p) })
}
func (s ProductService) Publish(ctx context.Context, id string, version int) error {
	p, e := s.Store.GetProduct(ctx, id)
	if e != nil {
		return e
	}
	if p.Version != version || p.Status != domain.ProductPending {
		return apperrors.ErrConflict
	}
	now := s.Clock.Now()
	rel := domain.ProductRelease{ID: now.Format("20060102150405.000000000"), ProductID: id, Version: version + 1, Status: "published", PublishedAt: &now, CreatedAt: now}
	return s.Store.WithTx(ctx, func(tx *sql.Tx) error { return s.Store.PublishProductTx(ctx, tx, id, version, rel) })
}
