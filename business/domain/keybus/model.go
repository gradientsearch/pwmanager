package keybus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/money"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/quantity"
)

// Key represents an individual key.
type Key struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
}

// NewKey is what we require from clients when adding a Key.
type NewKey struct {
	UserID   uuid.UUID
	Name     name.Name
	Cost     money.Money
	Quantity quantity.Quantity
}

// UpdateKey defines what information may be provided to modify an
// existing Key. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields so we can differentiate
// between a field that was not provided and a field that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types but
// we make exceptions around marshalling/unmarshalling.
type UpdateKey struct {
	Name     *name.Name
	Cost     *money.Money
	Quantity *quantity.Quantity
}
