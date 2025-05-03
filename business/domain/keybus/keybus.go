// Package keybus provides business access to key domain.
package keybus

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
	ErrNotFound     = errors.New("key not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, k Key) error
	Update(ctx context.Context, k Key) error
	Delete(ctx context.Context, k Key) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Key, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, keyID uuid.UUID) (Key, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Key, error)
	QueryByUserIDBundleID(ctx context.Context, bundleID uuid.UUID, userID uuid.UUID) (Key, error)
}

// Business manages the set of APIs for key access.
type Business struct {
	log      *logger.Logger
	userBus  *userbus.Business
	delegate *delegate.Delegate
	storer   Storer
}

// NewBusiness constructs a key business API for use.
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

// Create adds a new key to the system.
func (b *Business) Create(ctx context.Context, nk NewKey) (Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.create")
	defer span.End()

	usr, err := b.userBus.QueryByID(ctx, nk.UserID)
	if err != nil {
		return Key{}, fmt.Errorf("user.querybyid: %s: %w", nk.UserID, err)
	}

	if !usr.Enabled {
		return Key{}, ErrUserDisabled
	}

	now := time.Now()

	k := Key{
		ID:          uuid.New(),
		UserID:      nk.UserID,
		BundleID:    nk.BundleID,
		Data:        nk.Data,
		Roles:       nk.Roles,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, k); err != nil {
		return Key{}, fmt.Errorf("create: %w", err)
	}

	return k, nil
}

// Update modifies information about a key.
func (b *Business) Update(ctx context.Context, k Key, uk UpdateKey) (Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.update")
	defer span.End()

	if uk.Data != nil {
		k.Data = *uk.Data
	}

	k.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, k); err != nil {
		return Key{}, fmt.Errorf("update: %w", err)
	}

	return k, nil
}

// Delete removes the specified key.
func (b *Business) Delete(ctx context.Context, k Key) error {
	ctx, span := otel.AddSpan(ctx, "business.keybus.delete")
	defer span.End()

	if err := b.storer.Delete(ctx, k); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing keys.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.query")
	defer span.End()

	keys, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return keys, nil
}

// Count returns the total number of keys.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.count")
	defer span.End()

	return b.storer.Count(ctx, filter)
}

// QueryByID finds the key by the specified ID.
func (b *Business) QueryByID(ctx context.Context, keyID uuid.UUID) (Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.querybyid")
	defer span.End()

	k, err := b.storer.QueryByID(ctx, keyID)
	if err != nil {
		return Key{}, fmt.Errorf("query: keyID[%s]: %w", keyID, err)
	}

	return k, nil
}

// QueryByUserID finds the keys by a specified User ID.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.querybyuserid")
	defer span.End()

	keys, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return keys, nil
}

// QueryByUserIDBundleID finds the key by the specified userID and bundleID.
func (b *Business) QueryByUserIDBundleID(ctx context.Context, userID uuid.UUID, bundleID uuid.UUID) (Key, error) {
	ctx, span := otel.AddSpan(ctx, "business.keybus.querybybundleidanduserid")
	defer span.End()

	k, err := b.storer.QueryByUserIDBundleID(ctx, userID, bundleID)
	if err != nil {
		return Key{}, fmt.Errorf("query: userID[%s] bundleID[%s]: %w", userID, bundleID, err)
	}

	return k, nil
}
