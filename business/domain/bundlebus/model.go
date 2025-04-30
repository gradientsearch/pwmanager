package bundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

// Address represents an individual address.
type xAddress struct {
	Address1 string // We should create types for these fields.
	Address2 string
	ZipCode  string
	City     string
	State    string
	Country  string
}

// Bundle represents an individual bundle.
type Bundle struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Type        bundletype.BundleType
	DateCreated time.Time
	DateUpdated time.Time
}

// NewBundle is what we require from clients when adding a Bundle.
type NewBundle struct {
	UserID uuid.UUID
	Type   bundletype.BundleType
}

// UpdateBundle defines what information may be provided to modify an existing
// Bundle. All fields are optional so clients can send only the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank. Normally
// we do not want to use pointers to basic types but we make exception around
// marshalling/unmarshalling.
type UpdateBundle struct {
	Type *bundletype.BundleType
}
