package domain

import "time"

type Role string

const (
	RoleProvider  Role = "provider"
	RoleReviewer  Role = "reviewer"
	RoleDeveloper Role = "developer"
	RoleOperator  Role = "operator"
)

type ResourceStatus string

const (
	ResourceDraft      ResourceStatus = "draft"
	ResourceSubmitted  ResourceStatus = "submitted"
	ResourceReviewing  ResourceStatus = "under_review"
	ResourceRegistered ResourceStatus = "registered"
	ResourcePublished  ResourceStatus = "published"
	ResourceRejected   ResourceStatus = "rejected"
)

type AuthorizationStatus string

const (
	AuthRequested AuthorizationStatus = "requested"
	AuthApproved  AuthorizationStatus = "approved"
	AuthActive    AuthorizationStatus = "active"
	AuthRevoked   AuthorizationStatus = "revoked"
	AuthExpired   AuthorizationStatus = "expired"
)

type LeaseStatus string

const (
	LeaseReserved  LeaseStatus = "reserved"
	LeaseRunning   LeaseStatus = "running"
	LeaseCompleted LeaseStatus = "completed"
	LeaseFailed    LeaseStatus = "failed"
	LeaseExpired   LeaseStatus = "expired"
)

type ProductStatus string

const (
	ProductDraft     ProductStatus = "draft"
	ProductTesting   ProductStatus = "testing"
	ProductPending   ProductStatus = "pending_release"
	ProductPublished ProductStatus = "published"
	ProductWithdrawn ProductStatus = "withdrawn"
)

type Tenant struct {
	ID, Name  string
	CreatedAt time.Time
}
type User struct {
	ID, TenantID, Email string
	Role                Role
	Active              bool
	CreatedAt           time.Time
}
type Session struct {
	ID, UserID, TokenHash           string
	ExpiresAt, CreatedAt, RevokedAt *time.Time
}
type Resource struct {
	ID, TenantID, OwnerID, Code, Name, Description string
	Status                                         ResourceStatus
	Version                                        int
	CreatedAt, UpdatedAt                           time.Time
}
type ResourceVersion struct {
	ID, ResourceID      string
	Number              int
	Schema, ContentHash string
	CreatedAt           time.Time
}
type Review struct {
	ID, ResourceID, ReviewerID string
	Decision, Notes            string
	Status                     string
	LeaseUntil                 *time.Time
	CreatedAt                  time.Time
}
type Authorization struct {
	ID, TenantID, ResourceID, ApplicantID string
	Status                                AuthorizationStatus
	Purpose                               string
	Quota                                 int
	Used                                  int
	Version                               int
	CreatedAt, UpdatedAt                  time.Time
}
type Lease struct {
	ID, AuthorizationID, HolderID string
	Status                        LeaseStatus
	ExpiresAt                     time.Time
	Version                       int
	CreatedAt                     time.Time
}
type SandboxRun struct {
	ID, LeaseID           string
	Status                string
	Input, Output         string
	StartedAt, FinishedAt *time.Time
}
type Product struct {
	ID, TenantID, OwnerID, Name, ResourceID string
	Status                                  ProductStatus
	Version                                 int
	CreatedAt, UpdatedAt                    time.Time
}
type ProductRelease struct {
	ID, ProductID string
	Version       int
	Status        string
	PublishedAt   *time.Time
	CreatedAt     time.Time
}
type AuditEvent struct {
	ID, TenantID, ActorID, ObjectType, ObjectID, Action, Result, RequestID, Detail string
	CreatedAt                                                                      time.Time
}
type OutboxJob struct {
	ID, Topic, Payload string
	Attempts           int
	AvailableAt        time.Time
	LockedUntil        *time.Time
	Status             string
	LastError          string
	CreatedAt          time.Time
}
