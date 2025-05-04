// Package entrydb contains entry related CRUD functionality.
package entrydb

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for entry database access.
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
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (entrybus.Storer, error) {
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

// Create adds a Entry to the sqldb. It returns the created Entry with
// fields like ID and DateCreated populated.
func (s *Store) Create(ctx context.Context, k entrybus.Entry) error {
	const q = `
	INSERT INTO entries
		(entry_id, user_id, bundle_id, data, date_created, date_updated)
	VALUES
		(:entry_id, :user_id, :bundle_id, :data, :date_created, :date_updated)`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBEntry(k)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Update modifies data about a entrybus. It will error if the specified ID is
// invalid or does not reference an existing entrybus.
func (s *Store) Update(ctx context.Context, k entrybus.Entry) error {
	const q = `
	UPDATE
		entries
	SET
		"data" = :data,
		"user_id" = :user_id,
		"date_updated" = :date_updated
	WHERE
		entry_id = :entry_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, toDBEntry(k)); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Delete removes the entry identified by a given ID.
func (s *Store) Delete(ctx context.Context, k entrybus.Entry) error {
	data := struct {
		ID string `db:"entry_id"`
	}{
		ID: k.ID.String(),
	}

	const q = `
	DELETE FROM
		entries
	WHERE
		entry_id = :entry_id`

	if err := sqldb.NamedExecContext(ctx, s.log, s.db, q, data); err != nil {
		return fmt.Errorf("namedexeccontext: %w", err)
	}

	return nil
}

// Query gets all Entries from the database.
func (s *Store) Query(ctx context.Context, filter entrybus.QueryFilter, orderBy order.By, page page.Page) ([]entrybus.Entry, error) {
	data := map[string]any{
		"offset":        (page.Number() - 1) * page.RowsPerPage(),
		"rows_per_page": page.RowsPerPage(),
	}

	const q = `
	SELECT
	    entry_id, user_id,  bundle_id, data, date_created, date_updated
	FROM
		entries`

	buf := bytes.NewBufferString(q)
	s.applyFilter(filter, data, buf)

	orderByClause, err := orderByClause(orderBy)
	if err != nil {
		return nil, err
	}

	buf.WriteString(orderByClause)
	buf.WriteString(" OFFSET :offset ROWS FETCH NEXT :rows_per_page ROWS ONLY")

	var dbEntries []entry
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, buf.String(), data, &dbEntries); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	return toBusEntries(dbEntries)
}

// Count returns the total number of users in the DB.
func (s *Store) Count(ctx context.Context, filter entrybus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		entries`

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

// QueryByID finds the entry identified by a given ID.
func (s *Store) QueryByID(ctx context.Context, entryID uuid.UUID) (entrybus.Entry, error) {
	data := struct {
		ID string `db:"entry_id"`
	}{
		ID: entryID.String(),
	}

	const q = `
	SELECT
	    entry_id, user_id,  bundle_id, data, date_created, date_updated
	FROM
		entries
	WHERE
		entry_id = :entry_id`

	var dbEntry entry
	if err := sqldb.NamedQueryStruct(ctx, s.log, s.db, q, data, &dbEntry); err != nil {
		if errors.Is(err, sqldb.ErrDBNotFound) {
			return entrybus.Entry{}, fmt.Errorf("db: %w", entrybus.ErrNotFound)
		}
		return entrybus.Entry{}, fmt.Errorf("db: %w", err)
	}

	return toBusEntry(dbEntry)
}

// QueryByUserID finds the entry identified by a given User ID.
func (s *Store) QueryByUserID(ctx context.Context, userID uuid.UUID) ([]entrybus.Entry, error) {
	data := struct {
		ID string `db:"user_id"`
	}{
		ID: userID.String(),
	}

	const q = `
	SELECT
	    entry_id, user_id, date_created, date_updated
	FROM
		entries
	WHERE
		user_id = :user_id`

	var dbEntries []entry
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dbEntries); err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}

	return toBusEntries(dbEntries)
}
