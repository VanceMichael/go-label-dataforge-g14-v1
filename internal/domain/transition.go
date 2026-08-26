package domain

import "fmt"

type Transition struct {
	Entity string
	From   string
	To     string
	Actor  string
}

func (t Transition) Validate() error {
	if t.Entity == "" || t.From == "" || t.To == "" {
		return fmt.Errorf("transition fields are required")
	}
	switch t.Entity {
	case "resource":
		if !ResourceStatus(t.From).Can(ResourceStatus(t.To)) {
			return fmt.Errorf("resource transition denied")
		}
	case "authorization":
		if !AuthorizationStatus(t.From).Can(AuthorizationStatus(t.To)) {
			return fmt.Errorf("authorization transition denied")
		}
	case "lease":
		if !LeaseStatus(t.From).Can(LeaseStatus(t.To)) {
			return fmt.Errorf("lease transition denied")
		}
	case "product":
		if !ProductStatus(t.From).Can(ProductStatus(t.To)) {
			return fmt.Errorf("product transition denied")
		}
	default:
		return fmt.Errorf("unknown entity %s", t.Entity)
	}
	return nil
}
func ResourceLifecycle() []ResourceStatus {
	return []ResourceStatus{ResourceDraft, ResourceSubmitted, ResourceReviewing, ResourceRegistered, ResourcePublished, ResourceRejected}
}
func AuthorizationLifecycle() []AuthorizationStatus {
	return []AuthorizationStatus{AuthRequested, AuthApproved, AuthActive, AuthRevoked, AuthExpired}
}
func LeaseLifecycle() []LeaseStatus {
	return []LeaseStatus{LeaseReserved, LeaseRunning, LeaseCompleted, LeaseFailed, LeaseExpired}
}
func ProductLifecycle() []ProductStatus {
	return []ProductStatus{ProductDraft, ProductTesting, ProductPending, ProductPublished, ProductWithdrawn}
}
func IsTerminal(status string) bool {
	return status == string(ResourcePublished) || status == string(ResourceRejected) || status == string(AuthRevoked) || status == string(AuthExpired) || status == string(LeaseCompleted) || status == string(LeaseFailed) || status == string(LeaseExpired) || status == string(ProductPublished) || status == string(ProductWithdrawn)
}
