package domain

import "fmt"

func (r ResourceStatus) Can(next ResourceStatus) bool {
	switch r {
	case ResourceDraft:
		return next == ResourceSubmitted
	case ResourceSubmitted:
		return next == ResourceReviewing || next == ResourceRejected
	case ResourceReviewing:
		return next == ResourceRegistered || next == ResourceRejected
	case ResourceRegistered:
		return next == ResourcePublished
	case ResourcePublished:
		return next == ResourceRejected
	}
	return false
}
func (a AuthorizationStatus) Can(next AuthorizationStatus) bool {
	switch a {
	case AuthRequested:
		return next == AuthApproved
	case AuthApproved:
		return next == AuthActive || next == AuthRevoked
	case AuthActive:
		return next == AuthRevoked || next == AuthExpired
	}
	return false
}
func (l LeaseStatus) Can(next LeaseStatus) bool {
	switch l {
	case LeaseReserved:
		return next == LeaseRunning || next == LeaseExpired
	case LeaseRunning:
		return next == LeaseCompleted || next == LeaseFailed || next == LeaseExpired
	}
	return false
}
func (p ProductStatus) Can(next ProductStatus) bool {
	switch p {
	case ProductDraft:
		return next == ProductTesting
	case ProductTesting:
		return next == ProductPending || next == ProductDraft
	case ProductPending:
		return next == ProductPublished || next == ProductTesting
	case ProductPublished:
		return next == ProductWithdrawn
	}
	return false
}
func RequireTransition(from, to string, ok bool) error {
	if !ok {
		return fmt.Errorf("invalid transition %s -> %s", from, to)
	}
	return nil
}
