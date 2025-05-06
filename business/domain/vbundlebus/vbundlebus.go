// Package vbundlebus provides business access to view key domain.
package vbundlebus

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/foundation/otel"
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	QueryByID(ctx context.Context, userID uuid.UUID) ([]Key, error)
}

// Business manages the set of APIs for view key access.
type Business struct {
	storer Storer
}

// NewBusiness constructs a vbundle business API for use.
func NewBusiness(storer Storer) *Business {
	return &Business{
		storer: storer,
	}
}

// Query retrieves a list of existing keys.
func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) ([]Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.vbundlebus.query")
	defer span.End()

	users, err := b.storer.QueryByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}
