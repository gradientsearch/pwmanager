// Package entrybus provides business access to entry domain.
package entrybus

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/delegate"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/otel"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("entry not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, k Entry) error
	Update(ctx context.Context, k Entry) error
	Delete(ctx context.Context, k Entry) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Entry, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, entryID uuid.UUID) (Entry, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Entry, error)
}

// Business manages the set of APIs for entry access.
type Business struct {
	log      *logger.Logger
	userBus  *userbus.Business
	delegate *delegate.Delegate
	storer   Storer
}

// NewBusiness constructs a entry business API for use.
func NewBusiness(log *logger.Logger, userBus *userbus.Business, delegate *delegate.Delegate, storer Storer) *Business {
	b := Business{
		log:      log,
		userBus:  userBus,
		delegate: delegate,
		storer:   storer,
	}

	b.registerDelegateFunctions()

	return &b
}

// NewWithTx constructs a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (*Business, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bus := Business{
		log:      b.log,
		userBus:  userBus,
		delegate: b.delegate,
		storer:   storer,
	}

	return &bus, nil
}

// Create adds a new entry to the system.
func (b *Business) Create(ctx context.Context, ne NewEntry) (Entry, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.create")
	defer span.End()

	usr, err := b.userBus.QueryByID(ctx, ne.UserID)
	if err != nil {
		return Entry{}, fmt.Errorf("user.querybyid: %s: %w", ne.UserID, err)
	}

	if !usr.Enabled {
		return Entry{}, ErrUserDisabled
	}

	now := time.Now()

	e := Entry{
		ID:          uuid.New(),
		Data:        ne.Data,
		UserID:      ne.UserID,
		BundleID:    ne.BundleID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, e); err != nil {
		return Entry{}, fmt.Errorf("create: %w", err)
	}

	return e, nil
}

// Update modifies information about a entry.
func (b *Business) Update(ctx context.Context, e Entry, ue UpdateEntry) (Entry, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.update")
	defer span.End()

	if ue.Data != nil {
		e.Data = *ue.Data
	}

	e.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, e); err != nil {
		return Entry{}, fmt.Errorf("update: %w", err)
	}

	return e, nil
}

// Delete removes the specified entry.
func (b *Business) Delete(ctx context.Context, e Entry) error {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.delete")
	defer span.End()

	if err := b.storer.Delete(ctx, e); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing entries.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Entry, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.query")
	defer span.End()

	entries, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return entries, nil
}

// Count returns the total number of entries.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.count")
	defer span.End()

	return b.storer.Count(ctx, filter)
}

// QueryByID finds the entry by the specified ID.
func (b *Business) QueryByID(ctx context.Context, entryID uuid.UUID) (Entry, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.querybyid")
	defer span.End()

	k, err := b.storer.QueryByID(ctx, entryID)
	if err != nil {
		return Entry{}, fmt.Errorf("query: entryID[%s]: %w", entryID, err)
	}

	return k, nil
}

// QueryByUserID finds the entries by a specified User ID.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Entry, error) {
	ctx, span := otel.AddSpan(ctx, "business.entrybus.querybyuserid")
	defer span.End()

	entries, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return entries, nil
}
