// Package keydb contains key related CRUD functionality.
package keydb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for key database access.
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
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (keybus.Storer, error) {
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

// Create adds a Key to the sqldb. It returns the created Key with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, k keybus.Key) error {
	const q = `
	INSERT INTO keys
		(key_id, user_id, bundle_id, data, date_created, date_updated)
	VALUES
		(:key_id, :user_id, :bundle_id, :data, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBKey(k)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies data about a keybus. It will error if the specified ID is
// invalid or does not reference an existing keybus.
func (s *Store) Update(ctx context.Context, k keybus.Key) error {
	const q = `
	UPDATE
		keys
	SET
		"data" = :data,
		"date_updated" = :date_updated
	WHERE
		key_id = :key_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBKey(k)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes the key identified by a given ID.
func (s *Store) Delete(ctx context.Context, k keybus.Key) error {
	data := struct {
		ID string `db:"key_id"`
	}{
		ID: k.ID.String(),
	}

	const q = `
	DELETE FROM
		keys
	WHERE
		key_id = :key_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query gets all Keys from the database.
func (s *Store) Query(ctx context.Context, filter keybus.QueryFilter, orderBy order.By, page page.Page) ([]keybus.Key, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
	    key_id, user_id,  bundle_id, data, date_created, date_updated
	FROM
		keys`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbKeys []key
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbKeys); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusKeys(dbKeys)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter keybus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		keys`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	var count struct {
		Count   int `db:"count"`
		Sold    int `db:"sold"`
		Revenue int `db:"revenue"`
	}
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, buf.String(), data, &count); err != nil {
		return 0, fmt.Errorf("db: %w", err)
	}

	return count.Count, nil
}

// QueryByID finds the key identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, keyID uuid.UUID) (keybus.Key, error) {
	data := struct {
		ID string `db:"key_id"`
	}{
		ID: keyID.String(),
	}

	const q = `
	SELECT
	    key_id, user_id,  bundle_id, data, date_created, date_updated
	FROM
		keys
	WHERE
		key_id = :key_id`

	var dbKey key
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbKey); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return keybus.Key{}, fmt.Errorf("db: %w", keybus.ErrNotFound)
		}
		return keybus.Key{}, fmt.Errorf("db: %w", err)
	}

	return toBusKey(dbKey)
}

// QueryByUserID finds the key identified by a given User ID.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]keybus.Key, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    key_id, user_id, bundle_id, date_created, date_updated
	FROM
		keys
	WHERE
		user_id = :user_id`

	var dbKeys []key
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbKeys); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusKeys(dbKeys)
}

// QueryByUserIDBundleID finds the key identified by a given User ID and Bundle ID .
func (s *Store) QueryByUserIDBundleID(ctx context.Context, userID uuid.UUID, bundleID uuid.UUID) (keybus.Key, error) {
	data := struct {
		UserID   string `db:"user_id"`
		BundleID string `db:"bundle_id"`
	}{
		UserID:   userID.String(),
		BundleID: bundleID.String(),
	}

	const q = `
	SELECT
	    key_id, user_id, bundle_id, date_created, date_updated
	FROM
		keys
	WHERE
		user_id = :user_id 
		AND bundle_id = :bundle_id`

	var dbKey key
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbKey); err != nil {
		return keybus.Key{}, fmt.Errorf("db: %w", err)
	}

	return toBusKey(dbKey)
}
