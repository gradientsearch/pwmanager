package keydb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/types/money"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/quantity"
)

type key struct {
	ID          uuid.UUID `db:"key_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Cost        float64   `db:"cost"`
	Quantity    int       `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBKey(bus keybus.Key) key {
	db := key{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Name:        bus.Name.String(),
		Cost:        bus.Cost.Value(),
		Quantity:    bus.Quantity.Value(),
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return db
}

func toBusKey(db key) (keybus.Key, error) {
	name, err := name.Parse(db.Name)
	if err != nil {
		return keybus.Key{}, fmt.Errorf("parse name: %w", err)
	}

	cost, err := money.Parse(db.Cost)
	if err != nil {
		return keybus.Key{}, fmt.Errorf("parse cost: %w", err)
	}

	quantity, err := quantity.Parse(db.Quantity)
	if err != nil {
		return keybus.Key{}, fmt.Errorf("parse quantity: %w", err)
	}

	bus := keybus.Key{
		ID:          db.ID,
		UserID:      db.UserID,
		Name:        name,
		Cost:        cost,
		Quantity:    quantity,
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
