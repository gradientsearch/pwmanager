// Package bundlebus provides business access to bundle domain.
package bundlebus

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
	ErrNotFound     = errors.New("bundle not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behaviour this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, hme Bundle) error
	Update(ctx context.Context, hme Bundle) error
	Delete(ctx context.Context, hme Bundle) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Bundle, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, bundleID uuid.UUID) (Bundle, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Bundle, error)
}

// Business manages the set of APIs for bundle api access.
type Business struct {
	log      *logger.Logger
	userBus  *userbus.Business
	delegate *delegate.Delegate
	storer   Storer
}

// NewBusiness constructs a bundle business API for use.
func NewBusiness(log *logger.Logger, userBus *userbus.Business, delegate *delegate.Delegate, storer Storer) *Business {
	return &Business{
		log:      log,
		userBus:  userBus,
		delegate: delegate,
		storer:   storer,
	}
}

// NewWithTx constructs a new domain value that will use the
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

// Create adds a new bundle to the system.
func (b *Business) Create(ctx context.Context, nh NewBundle) (Bundle, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.create")
	defer span.End()

	usr, err := b.userBus.QueryByID(ctx, nh.UserID)
	if err != nil {
		return Bundle{}, fmt.Errorf("user.querybyid: %s: %w", nh.UserID, err)
	}

	if !usr.Enabled {
		return Bundle{}, ErrUserDisabled
	}

	now := time.Now()

	hme := Bundle{
		ID:   uuid.New(),
		Type: nh.Type,
		Address: Address{
			Address1: nh.Address.Address1,
			Address2: nh.Address.Address2,
			ZipCode:  nh.Address.ZipCode,
			City:     nh.Address.City,
			State:    nh.Address.State,
			Country:  nh.Address.Country,
		},
		UserID:      nh.UserID,
		DateCreated: now,
		DateUpdated: now,
	}

	if err := b.storer.Create(ctx, hme); err != nil {
		return Bundle{}, fmt.Errorf("create: %w", err)
	}

	return hme, nil
}

// Update modifies information about a bundle.
func (b *Business) Update(ctx context.Context, hme Bundle, uh UpdateBundle) (Bundle, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.update")
	defer span.End()

	if uh.Type != nil {
		hme.Type = *uh.Type
	}

	if uh.Address != nil {
		if uh.Address.Address1 != nil {
			hme.Address.Address1 = *uh.Address.Address1
		}

		if uh.Address.Address2 != nil {
			hme.Address.Address2 = *uh.Address.Address2
		}

		if uh.Address.ZipCode != nil {
			hme.Address.ZipCode = *uh.Address.ZipCode
		}

		if uh.Address.City != nil {
			hme.Address.City = *uh.Address.City
		}

		if uh.Address.State != nil {
			hme.Address.State = *uh.Address.State
		}

		if uh.Address.Country != nil {
			hme.Address.Country = *uh.Address.Country
		}
	}

	hme.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, hme); err != nil {
		return Bundle{}, fmt.Errorf("update: %w", err)
	}

	return hme, nil
}

// Delete removes the specified bundle.
func (b *Business) Delete(ctx context.Context, hme Bundle) error {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.delete")
	defer span.End()

	if err := b.storer.Delete(ctx, hme); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Query retrieves a list of existing bundles.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Bundle, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.query")
	defer span.End()

	hmes, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}

// Count returns the total number of bundles.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.count")
	defer span.End()

	return b.storer.Count(ctx, filter)
}

// QueryByID finds the bundle by the specified ID.
func (b *Business) QueryByID(ctx context.Context, bundleID uuid.UUID) (Bundle, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.querybyid")
	defer span.End()

	hme, err := b.storer.QueryByID(ctx, bundleID)
	if err != nil {
		return Bundle{}, fmt.Errorf("query: bundleID[%s]: %w", bundleID, err)
	}

	return hme, nil
}

// QueryByUserID finds the bundles by a specified User ID.
func (b *Business) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Bundle, error) {
	ctx, span := otel.AddSpan(ctx, "business.bundlebus.querybyuserid")
	defer span.End()

	hmes, err := b.storer.QueryByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return hmes, nil
}
