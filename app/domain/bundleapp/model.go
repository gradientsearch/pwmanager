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

// Bundle represents information about an individual bundle.
type Bundle struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Type        string `json:"type"`
	Metadata    string `json:"metadata"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Bundle) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppBundle(b bundlebus.Bundle) Bundle {
	return Bundle{
		ID:       b.ID.String(),
		UserID:   b.UserID.String(),
		Type:     b.Type.String(),
		Metadata: b.Metadata,

		DateCreated: b.DateCreated.Format(time.RFC3339),
		DateUpdated: b.DateUpdated.Format(time.RFC3339),
	}
}

func toAppBundles(bundles []bundlebus.Bundle) []Bundle {
	app := make([]Bundle, len(bundles))
	for i, b := range bundles {
		app[i] = toAppBundle(b)
	}

	return app
}

// =============================================================================

// NewBundle defines the data needed to add a new bundle.
type NewBundle struct {
	Type     string `json:"type" validate:"required"`
	Metadata string `json:"metadata" validate:"required"`
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
		UserID:   userID,
		Type:     typ,
		Metadata: app.Metadata,
	}

	return bus, nil
}

// =============================================================================

// UpdateBundle defines the data needed to update a bundle.
type UpdateBundle struct {
	Type     *string `json:"type"` // TODO may not want to allow updating bundle type 🤔
	Metadata *string `json:"metadata"`
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
		Type:     &t,
		Metadata: app.Metadata,
	}

	return bus, nil
}
