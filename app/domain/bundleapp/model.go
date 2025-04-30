package bundleapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

// Address represents information about an individual address.
type Address struct {
	Address1 string `json:"address1"`
	Address2 string `json:"address2"`
	ZipCode  string `json:"zipCode"`
	City     string `json:"city"`
	State    string `json:"state"`
	Country  string `json:"country"`
}

// Bundle represents information about an individual bundle.
type Bundle struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userID"`
	Type        string  `json:"type"`
	Address     Address `json:"address"`
	DateCreated string  `json:"dateCreated"`
	DateUpdated string  `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Bundle) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppBundle(hme bundlebus.Bundle) Bundle {
	return Bundle{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type.String(),

		DateCreated: hme.DateCreated.Format(time.RFC3339),
		DateUpdated: hme.DateUpdated.Format(time.RFC3339),
	}
}

func toAppBundles(bundles []bundlebus.Bundle) []Bundle {
	app := make([]Bundle, len(bundles))
	for i, hme := range bundles {
		app[i] = toAppBundle(hme)
	}

	return app
}

// =============================================================================

// NewBundle defines the data needed to add a new bundle.
type NewBundle struct {
	Type string `json:"type" validate:"required"`
}

// Decode implements the decoder interface.
func (app *NewBundle) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks if the data in the model is considered clean.
func (app NewBundle) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewBundle(ctx context.Context, app NewBundle) (bundlebus.NewBundle, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return bundlebus.NewBundle{}, fmt.Errorf("getuserid: %w", err)
	}

	typ, err := bundletype.Parse(app.Type)
	if err != nil {
		return bundlebus.NewBundle{}, fmt.Errorf("parse: %w", err)
	}

	bus := bundlebus.NewBundle{
		UserID: userID,
		Type:   typ,
	}

	return bus, nil
}

// =============================================================================

// UpdateAddress defines the data needed to update an address.
type UpdateAddress struct {
	Address1 *string `json:"address1" validate:"omitempty,min=1,max=70"`
	Address2 *string `json:"address2" validate:"omitempty,max=70"`
	ZipCode  *string `json:"zipCode" validate:"omitempty,numeric"`
	City     *string `json:"city"`
	State    *string `json:"state" validate:"omitempty,min=1,max=48"`
	Country  *string `json:"country" validate:"omitempty,iso3166_1_alpha2"`
}

// UpdateBundle defines the data needed to update a bundle.
type UpdateBundle struct {
	Type    *string        `json:"type"`
	Address *UpdateAddress `json:"address"`
}

// Decode implements the decoder interface.
func (app *UpdateBundle) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateBundle) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateBundle(app UpdateBundle) (bundlebus.UpdateBundle, error) {
	var t bundletype.BundleType
	if app.Type != nil {
		var err error
		t, err = bundletype.Parse(*app.Type)
		if err != nil {
			return bundlebus.UpdateBundle{}, fmt.Errorf("parse: %w", err)
		}
	}

	bus := bundlebus.UpdateBundle{
		Type: &t,
	}

	return bus, nil
}
