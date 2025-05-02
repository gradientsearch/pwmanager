package entryapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/types/entry"
)

// Entry represents information about an individual entry.
type Entry struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	BundleID    string `json:"bundleID"`
	Data        string `json:"data"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Entry) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppEntry(e entrybus.Entry) Entry {
	return Entry{
		ID:          e.ID.String(),
		BundleID:    e.BundleID.String(),
		UserID:      e.UserID.String(),
		Data:        e.Data.String(),
		DateCreated: e.DateCreated.Format(time.RFC3339),
		DateUpdated: e.DateUpdated.Format(time.RFC3339),
	}
}

func toAppEntries(entries []entrybus.Entry) []Entry {
	app := make([]Entry, len(entries))
	for i, k := range entries {
		app[i] = toAppEntry(k)
	}

	return app
}

// =============================================================================

// NewEntry defines the data needed to add a new entry.
type NewEntry struct {
	Data string `json:"data" validate:"required"`
}

// Decode implements the decoder interface.
func (app *NewEntry) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewEntry) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewEntry(ctx context.Context, app NewEntry) (entrybus.NewEntry, error) {
	ne, err := mid.GetEntry(ctx)
	if err != nil {
		return entrybus.NewEntry{}, fmt.Errorf("getentry: %w", err)
	}

	data, err := entry.Parse(app.Data)
	if err != nil {
		return entrybus.NewEntry{}, fmt.Errorf("parse data: %w", err)
	}

	bus := entrybus.NewEntry{
		UserID:   ne.UserID,
		BundleID: ne.BundleID,
		Data:     data,
	}

	return bus, nil
}

// =============================================================================

// UpdateEntry defines the data needed to update a entry.
type UpdateEntry struct {
	Data *string `json:"data"`
}

// Decode implements the decoder interface.
func (app *UpdateEntry) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app UpdateEntry) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusUpdateEntry(app UpdateEntry) (entrybus.UpdateEntry, error) {
	var kd *entry.Entry
	if app.Data != nil {
		k, err := entry.Parse(*app.Data)
		if err != nil {
			return entrybus.UpdateEntry{}, fmt.Errorf("parse: %w", err)
		}
		kd = &k
	}

	bus := entrybus.UpdateEntry{
		Data: kd,
	}

	return bus, nil
}
