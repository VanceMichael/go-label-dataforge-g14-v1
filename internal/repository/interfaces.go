package repository

import (
	"context"
	"database/sql"
	"github.com/VanceMichael/go-label-dataforge-g14-v1/internal/domain"
	"time"
)

type Store interface {
	WithTx(context.Context, func(*sql.Tx) error) error
	CreateTenant(context.Context, *sql.Tx, domain.Tenant) error
	CreateUser(context.Context, *sql.Tx, domain.User) error
	FindUserByEmail(context.Context, string, string) (domain.User, error)
	CreateSession(context.Context, domain.Session) error
	FindSession(context.Context, string) (domain.Session, error)
	RevokeSession(context.Context, string) error
	CreateResourceTx(context.Context, *sql.Tx, domain.Resource, domain.ResourceVersion, domain.Review) error
	GetResource(context.Context, string) (domain.Resource, error)
	UpdateResourceStatusTx(context.Context, *sql.Tx, string, domain.ResourceStatus, int) error
	CreateAuthorizationTx(context.Context, *sql.Tx, domain.Authorization) error
	GetAuthorization(context.Context, string) (domain.Authorization, error)
	ApproveAuthorizationTx(context.Context, *sql.Tx, string, int) error
	CreateLeaseTx(context.Context, *sql.Tx, domain.Lease) error
	GetLease(context.Context, string) (domain.Lease, error)
	UpdateLeaseTx(context.Context, *sql.Tx, string, domain.LeaseStatus, int) error
	CreateRunTx(context.Context, *sql.Tx, domain.SandboxRun) error
	CreateProductTx(context.Context, *sql.Tx, domain.Product) error
	GetProduct(context.Context, string) (domain.Product, error)
	PublishProductTx(context.Context, *sql.Tx, string, int, domain.ProductRelease) error
	AddAuditTx(context.Context, *sql.Tx, domain.AuditEvent) error
	EnqueueTx(context.Context, *sql.Tx, domain.OutboxJob) error
	ClaimOutbox(context.Context, time.Time) (domain.OutboxJob, error)
	FinishOutbox(context.Context, string, error) error
	ExpireLeases(context.Context, time.Time) (int, error)
	ExpireReviews(context.Context, time.Time) (int, error)
}
