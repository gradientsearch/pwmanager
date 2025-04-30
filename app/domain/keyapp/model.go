package keyapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/types/money"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/quantity"
)

// Key represents information about an individual key.
type Key struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Name        string  `json:"name"`
	Cost        float64 `json:"cost"`
	Quantity    int     `json:"quantity"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Key) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppKey(prd keybus.Key) Key {
	return Key{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppKeys(prds []keybus.Key) []Key {
	app := make([]Key, len(prds))
	for i, prd := range prds {
		app[i] = toAppKey(prd)
	}

	return app
}

// =============================================================================

// NewKey defines the data needed to add a new key.
type NewKey struct {
	Name     string  `json:"name" validate:"required"`
	Cost     float64 `json:"cost" validate:"required,gte=0"`
	Quantity int     `json:"quantity" validate:"required,gte=1"`
}

// Decode implements the decoder interface.
func (app *NewKey) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewKey) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewKey(ctx context.Context, app NewKey) (keybus.NewKey, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("getuserid: %w", err)
	}

	name, err := name.Parse(app.Name)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse name: %w", err)
	}

	cost, err := money.Parse(app.Cost)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse cost: %w", err)
	}

	quantity, err := quantity.Parse(app.Quantity)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse quantity: %w", err)
	}

	bus := keybus.NewKey{
		UserID:   userID,
		Name:     name,
		Cost:     cost,
		Quantity: quantity,
	}

	return bus, nil
}

// =============================================================================

// UpdateKey defines the data needed to update a key.
type UpdateKey struct {
	Name     *string  `json:"name"`
	Cost     *float64 `json:"cost" validate:"omitempty,gte=0"`
	Quantity *int     `json:"quantity" validate:"omitempty,gte=1"`
}

// Decode implements the decoder interface.
func (app *UpdateKey) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateKey) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateKey(app UpdateKey) (keybus.UpdateKey, error) {
	var nme *name.Name
	if app.Name != nil {
		nm, err := name.Parse(*app.Name)
		if err != nil {
			return keybus.UpdateKey{}, fmt.Errorf("parse: %w", err)
		}
		nme = &nm
	}

	var cost *money.Money
	if app.Cost != nil {
		cst, err := money.Parse(*app.Cost)
		if err != nil {
			return keybus.UpdateKey{}, fmt.Errorf("parse: %w", err)
		}
		cost = &cst
	}

	var qnt *quantity.Quantity
	if app.Cost != nil {
		qn, err := quantity.Parse(*app.Quantity)
		if err != nil {
			return keybus.UpdateKey{}, fmt.Errorf("parse: %w", err)
		}
		qnt = &qn
	}

	bus := keybus.UpdateKey{
		Name:     nme,
		Cost:     cost,
		Quantity: qnt,
	}

	return bus, nil
}
