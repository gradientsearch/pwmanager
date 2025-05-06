// Package vbundledb provides access to the key view.
package vbundledb

import (
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
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) ([]vbundlebus.UserBundleKey, error) {
	data := map[string]any{
		"user_id": userID,
	}

	const q = `
SELECT
	u.user_id,
	u.name,
    b.bundle_id,
    b.type,
    b.metadata,
    b.date_created,
    b.date_updated,
    k.data AS key_data,
    k.roles AS key_roles,
    (
        SELECT json_agg(json_build_object('user_id', ku.user_id, 'name', ku.name, 'email', ku.email, 'roles', k2.roles))
        FROM keys k2
        JOIN users ku ON k2.user_id = ku.user_id
        WHERE k2.bundle_id = b.bundle_id
    ) AS users
FROM
    users u
JOIN
    bundles b ON b.user_id = u.user_id
LEFT JOIN
    keys k ON k.bundle_id = b.bundle_id AND k.user_id = u.user_id
WHERE b.user_id = :user_id`

	var dnKey []userBundleKey
	if err := sqldb.NamedQuerySlice(ctx, s.log, s.db, q, data, &dnKey); err != nil {
		return nil, fmt.Errorf("namedqueryslice: %w", err)
	}

	k, err := toBusBundles(dnKey)
	if err != nil {
		return nil, err
	}

	return k, nil
}
