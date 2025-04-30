package keydb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	kt "github.com/gradientsearch/pwmanager/business/types/key"
)

type key struct {
	ID          uuid.UUID `db:"key_id"`
	UserID      uuid.UUID `db:"user_id"`
	BundleID    uuid.UUID `db:"bundle_id"`
	Data        string    `db:"data"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBKey(bus keybus.Key) key {
	db := key{
		ID:          bus.ID,
		UserID:      bus.UserID,
		BundleID:    bus.BundleID,
		Data:        bus.Data.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return db
}

func toBusKey(db key) (keybus.Key, error) {
	key, err := kt.Parse(db.Data)
	if err != nil {
		return keybus.Key{}, fmt.Errorf("parse key: %w", err)
	}

	bus := keybus.Key{
		ID:          db.ID,
		UserID:      db.UserID,
		BundleID:    db.BundleID,
		Data:        key,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusKeys(dbs []key) ([]keybus.Key, error) {
	bus := make([]keybus.Key, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusKey(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
