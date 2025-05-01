package vbundledb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	kt "github.com/gradientsearch/pwmanager/business/types/key"
)

type key struct {
	ID          uuid.UUID `db:"key_id"`
	UserID      uuid.UUID `db:"user_id"`
	Name        string    `db:"name"`
	Cost        float64   `db:"cost"`
	Quantity    int       `db:"quantity"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
	UserName    string    `db:"user_name"`
}

func toBusKey(db key) (vbundlebus.Key, error) {
	name, err := kt.Parse(db.Name)
	if err != nil {
		return vbundlebus.Key{}, fmt.Errorf("parse name: %w", err)
	}

	bus := vbundlebus.Key{
		ID:     db.ID,
		UserID: db.UserID,
		Data:   name,

		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusKeys(dbKeys []key) ([]vbundlebus.Key, error) {
	bus := make([]vbundlebus.Key, len(dbKeys))

	for i, dbKey := range dbKeys {
		var err error
		bus[i], err = toBusKey(dbKey)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}
