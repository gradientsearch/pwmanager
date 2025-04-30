package bundledb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

type bundle struct {
	ID          uuid.UUID `db:"bundle_id"`
	UserID      uuid.UUID `db:"user_id"`
	Type        string    `db:"type"`
	Address1    string    `db:"address_1"`
	Address2    string    `db:"address_2"`
	ZipCode     string    `db:"zip_code"`
	City        string    `db:"city"`
	Country     string    `db:"country"`
	State       string    `db:"state"`
	DateCreated time.Time `db:"date_created"`
	DateUpdated time.Time `db:"date_updated"`
}

func toDBBundle(bus bundlebus.Bundle) bundle {
	db := bundle{
		ID:          bus.ID,
		UserID:      bus.UserID,
		Type:        bus.Type.String(),
		Address1:    bus.Address.Address1,
		Address2:    bus.Address.Address2,
		ZipCode:     bus.Address.ZipCode,
		City:        bus.Address.City,
		Country:     bus.Address.Country,
		State:       bus.Address.State,
		DateCreated: bus.DateCreated.UTC(),
		DateUpdated: bus.DateUpdated.UTC(),
	}

	return db
}

func toBusBundle(db bundle) (bundlebus.Bundle, error) {
	typ, err := bundletype.Parse(db.Type)
	if err != nil {
		return bundlebus.Bundle{}, fmt.Errorf("parse type: %w", err)
	}

	bus := bundlebus.Bundle{
		ID:     db.ID,
		UserID: db.UserID,
		Type:   typ,
		Address: bundlebus.Address{
			Address1: db.Address1,
			Address2: db.Address2,
			ZipCode:  db.ZipCode,
			City:     db.City,
			Country:  db.Country,
			State:    db.State,
		},
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusBundles(dbs []bundle) ([]bundlebus.Bundle, error) {
	bus := make([]bundlebus.Bundle, len(dbs))

	for i, db := range dbs {
		var err error
		bus[i], err = toBusBundle(db)
		if err != nil {
			return nil, fmt.Errorf("parse type: %w", err)
		}
	}

	return bus, nil
}
