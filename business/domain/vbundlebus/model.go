package vbundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/key"
)

// Key represents an individual key with extended information.
type Key struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Data        key.Key
	DateCreated time.Time
	DateUpdated time.Time
}
