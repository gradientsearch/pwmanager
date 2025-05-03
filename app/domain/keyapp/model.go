package keyapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/key"
)

// Key represents information about an individual key.
type Key struct {
	ID          string   `json:"id"`
	UserID      string   `json:"userID"`
	BundleID    string   `json:"bundleID"`
	Data        string   `json:"data"`
	Roles       []string `json:"roles"`
	DateCreated string   `json:"dateCreated"`
	DateUpdated string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Key) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppKey(k keybus.Key) Key {
	return Key{
		ID:          k.ID.String(),
		UserID:      k.UserID.String(),
		BundleID:    k.BundleID.String(),
		Data:        k.Data.String(),
		Roles:       bundlerole.ParseToString(k.Roles),
		DateCreated: k.DateCreated.Format(time.RFC3339),
		DateUpdated: k.DateUpdated.Format(time.RFC3339),
	}
}

func toAppKeys(keys []keybus.Key) []Key {
	app := make([]Key, len(keys))
	for i, k := range keys {
		app[i] = toAppKey(k)
	}

	return app
}

// =============================================================================

// NewKey defines the data needed to add a new key.
type NewKey struct {
	Data     string   `json:"data" validate:"required"`
	BundleID string   `json:"bundleID" validate:"required"`
	UserID   string   `json:"userID" validate:"required"`
	Roles    []string `json:"roles" validate:"required"`
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

	bundleID, err := uuid.Parse(app.BundleID)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("getuserid: %w", err)
	}

	data, err := key.Parse(app.Data)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse data: %w", err)
	}

	roles, err := bundlerole.ParseMany(app.Roles)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse: %w", err)
	}

	bus := keybus.NewKey{
		UserID:   userID,
		BundleID: bundleID,
		Roles:    roles,
		Data:     data,
	}

	return bus, nil
}

// =============================================================================

// UpdateKey defines the data needed to update a key.
type UpdateKey struct {
	Data *string `json:"data"`
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
	var kd *key.Key
	if app.Data != nil {
		k, err := key.Parse(*app.Data)
		if err != nil {
			return keybus.UpdateKey{}, fmt.Errorf("parse: %w", err)
		}
		kd = &k
	}

	bus := keybus.UpdateKey{
		Data: kd,
	}

	return bus, nil
}

// =============================================================================

// UpdateBundleRole defines the data needed to update a user role.
type UpdateBundleRole struct {
	Roles []string `json:"roles" validate:"required"`
}

// Decode implements the decoder interface.
func (app *UpdateBundleRole) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateBundleRole) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateBundleRole(app UpdateBundleRole) (keybus.UpdateKey, error) {
	var roles []bundlerole.Role
	if app.Roles != nil {
		var err error
		roles, err = bundlerole.ParseMany(app.Roles)
		if err != nil {
			return keybus.UpdateKey{}, fmt.Errorf("parse: %w", err)
		}
	}

	bus := keybus.UpdateKey{
		Roles: roles,
	}

	return bus, nil
}
