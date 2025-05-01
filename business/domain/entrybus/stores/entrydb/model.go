package entrydb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	kt "github.com/gradientsearch/pwmanager/business/types/entry"
)

type entry struct {
	ID          uuid.UUID `db:"entry_id"`
	UserID      uuid.UUID `db:"user_id"`
	BundleID    uuid.UUID `db:"bundle_id"`
	Data        string    `db:"data"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBEntry(bus entrybus.Entry) entry {
	db := entry{
		ID:          bus.ID,
		UserID:      bus.UserID,
		BundleID:    bus.BundleID,
		Data:        bus.Data.String(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return db
}

func toBusEntry(db entry) (entrybus.Entry, error) {
	entry, err := kt.Parse(db.Data)
	if err != nil {
		return entrybus.Entry{}, fmt.Errorf("parse entry: %w", err)
	}

	bus := entrybus.Entry{
		ID:          db.ID,
		UserID:      db.UserID,
		BundleID:    db.BundleID,
		Data:        entry,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusEntries(dbs []entry) ([]entrybus.Entry, error) {
	bus := make([]entrybus.Entry, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusEntry(db)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
