package keybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/key"
)

// Key represents an individual key.
type Key struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	BundleID    uuid.UUID
	Data        key.Key
	DateCreated time.Time
	DateUpdated time.Time
}

// NewKey is what we require from clients when adding a Key.
type NewKey struct {
	UserID   uuid.UUID
	BundleID uuid.UUID
	Data     key.Key
}

// UpdateKey defines what information may be provided to modify an
// existing Key. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateKey struct {
	Data   *key.Key
	UserID *uuid.UUID
}
