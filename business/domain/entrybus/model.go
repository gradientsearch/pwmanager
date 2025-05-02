package entrybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/entry"
)

// Entry represents an individual entry.
type Entry struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	BundleID    uuid.UUID
	Data        entry.Entry
	DateCreated time.Time
	DateUpdated time.Time
}

// NewEntry is what we require from clients when adding a Entry.
type NewEntry struct {
	UserID   uuid.UUID
	BundleID uuid.UUID
	Data     entry.Entry
}

// UpdateEntry defines what information may be provided to modify an
// existing Entry. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateEntry struct {
	Data   *entry.Entry
	UserID *uuid.UUID
}
