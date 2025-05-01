package bundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

// Bundle represents an individual bundle.
type Bundle struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        bundletype.BundleType
	Metadata    string
	DateCreated time.Time
	DateUpdated time.Time
}

// NewBundle is what we require from clients when adding a Bundle.
type NewBundle struct {
	UserID   uuid.UUID
	Type     bundletype.BundleType
	Metadata string
}

// UpdateBundle defines what information may be provided to modify an existing
// Bundle. All fields are optional so clients can send only the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exception around
// marshalling/unmarshalling.
type UpdateBundle struct {
	Type     *bundletype.BundleType
	Metadata *string
}
