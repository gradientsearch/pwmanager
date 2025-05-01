// Package bundledb contains bundle related CRUD functionality.
package bundledb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for bundle database access.
type Store struct {
	log *logger.Logger
	db  sqlx.ExtContext
}

// NewStore constructs the api for data access.
func NewStore(log *logger.Logger, db *sqlx.DB) *Store {
	return &Store{
		log: log,
		db:  db,
	}
}

// NewWithTx constructs a new Store value replacing the sqlx DB
// value with a sqlx DB value that is currently inside a transaction.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (bundlebus.Storer, error) {
	ec, err := sqldb.GetExtContext(tx)
	if err != nil {
		return nil, err
	}

	store := Store{
		log: s.log,
		db:  ec,
	}

	return &store, nil
}

// Create inserts a new bundle into the database.
func (s *Store) Create(ctx context.Context, bdl bundlebus.Bundle) error {
	const q = `
    INSERT INTO bundles
        (bundle_id, user_id, type, metadata, date_created, date_updated)
    VALUES
        (:bundle_id, :user_id, :type, :metadata, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBBundle(bdl)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes a bundle from the database.
func (s *Store) Delete(ctx context.Context, bdl bundlebus.Bundle) error {
	data := struct {
		ID string `db:"bundle_id"`
	}{
		ID: bdl.ID.String(),
	}

	const q = `
    DELETE FROM
	    bundles
	WHERE
	  	bundle_id = :bundle_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update replaces a bundle document in the database.
func (s *Store) Update(ctx context.Context, bdl bundlebus.Bundle) error {
	const q = `
    UPDATE
        bundles
    SET
        "type"          = :type,
		"metadata"      = :metadata,
        "date_updated"  = :date_updated
    WHERE
        bundle_id = :bundle_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBBundle(bdl)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query retrieves a list of existing bundles from the database.
func (s *Store) Query(ctx context.Context, filter bundlebus.QueryFilter, orderBy order.By, page page.Page) ([]bundlebus.Bundle, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
    SELECT
	    bundle_id, user_id, type, metadata, date_created, date_updated
	FROM
	  	bundles`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbBdls []bundle
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbBdls); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	bdls, err := toBusBundles(dbBdls)
	if err != nil {
		return nil, err
	}

	return bdls, nil
}

// Count returns the total number of bundles in the DB.
func (s *Store) Count(ctx context.Context, filter bundlebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
    SELECT
        count(1)
    FROM
        bundles`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count int `db:"count"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID gets the specified bundle from the database.
func (s *Store) QueryByID(ctx context.Context, bundleID uuid.UUID) (bundlebus.Bundle, error) {
	data := struct {
		ID string `db:"bundle_id"`
	}{
		ID: bundleID.String(),
	}

	const q = `
    SELECT
	  	bundle_id, user_id, type, metadata, date_created, date_updated
    FROM
        bundles
    WHERE
        bundle_id = :bundle_id`

	var dbBdl bundle
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbBdl); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return bundlebus.Bundle{}, fmt.Errorf("db: %w", bundlebus.ErrNotFound)
		}
		return bundlebus.Bundle{}, fmt.Errorf("db: %w", err)
	}

	return toBusBundle(dbBdl)
}

// QueryByUserID gets the specified bundle from the database by user id.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]bundlebus.Bundle, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    bundle_id, user_id, type, metadata, date_created, date_updated
	FROM
		bundles
	WHERE
		user_id = :user_id`

	var dbBdls []bundle
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbBdls); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusBundles(dbBdls)
}
