// Package vbundledb provides access to the key view.
package vbundledb

import (
	"bytes"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// Store manages the set of APIs for key view database access.
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

// Query retrieves a list of existing keys from the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) ([]vbundlebus.Key, error) {
	data := map[string]any{
		"user_id": userID,
	}

	const q = `
SELECT
	u.user_id,
    b.bundle_id,
    b.type AS bundle_type,
    b.metadata AS bundle_metadata,
    b.date_created AS bundle_date_created,
    b.date_updated AS bundle_date_updated,
    k.key_id,
    k.data AS key_data,
    k.roles AS key_roles,
    k.date_created AS key_date_created,
    k.date_updated AS key_date_updated,
    (
        SELECT json_agg(json_build_object('user_id', ku.user_id, 'name', ku.name, 'email', ku.email, 'roles', k2.roles))
        FROM keys k2
        JOIN users ku ON k2.user_id = ku.user_id
        WHERE k2.bundle_id = b.bundle_id
    ) AS users_with_access
FROM
    users u
JOIN
    bundles b ON b.user_id = u.user_id
LEFT JOIN
    keys k ON k.bundle_id = b.bundle_id AND k.user_id = u.user_id
WHERE b.user_id = :user_id`

	var dnKey []key
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dnKey); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	k, err := toBusKeys(dnKey)
	if err != nil {
		return nil, err
	}

	return k, nil
}

// Count returns the total number of keys in the DB.
func (s *Store) Count(ctx context.Context, filter vbundlebus.QueryFilter) (int, error) {
	data := map[string]any{}

	const q = `
	SELECT
		count(1)
	FROM
		view_keys`

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
