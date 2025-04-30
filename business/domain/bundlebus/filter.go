package bundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

// QueryFilter holds the available fields a query can be filtered on.
// We are using pointer semantics because the With API mutates the value.
type QueryFilter struct {
	ID               *uuid.UUID
	UserID           *uuid.UUID
	Type             *bundletype.BundleType
	StartCreatedDate *time.Time
	EndCreatedDate   *time.Time
}
