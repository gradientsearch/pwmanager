package vbundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/money"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/quantity"
)

// Product represents an individual product with extended information.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	Cost        money.Money
	Quantity    quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
	UserName    name.Name
}
