package entryapp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/types/entry"
)

// Entry represents information about an individual entry.
type EntryTx struct {
	Entry  Entry  `json:"entry"`
	Bundle Bundle `json:"bundle"`
}

// Encode implements the encoder interface.
func (app EntryTx) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppEntryTx(e entrybus.Entry, b bundlebus.Bundle) EntryTx {
	return EntryTx{
		Entry: Entry{
			ID:          e.ID.String(),
			BundleID:    e.BundleID.String(),
			UserID:      e.UserID.String(),
			Data:        e.Data.String(),
			DateCreated: e.DateCreated.Format(time.RFC3339),
			DateUpdated: e.DateUpdated.Format(time.RFC3339),
		},
		Bundle: Bundle{
			ID:       b.ID.String(),
			UserID:   b.UserID.String(),
			Type:     b.Type.String(),
			Metadata: b.Metadata,

			DateCreated: b.DateCreated.Format(time.RFC3339),
			DateUpdated: b.DateUpdated.Format(time.RFC3339),
		},
	}
}

// ------

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

// NewEntryTX defines the data needed to add a new entry.
type NewEntryTX struct {
	Data     string `json:"data" validate:"required"`
	Metadata string `json:"metadatadata" validate:"required"`
}

// Decode implements the decoder interface.
func (app *NewEntryTX) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app NewEntryTX) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewEntry(ctx context.Context, app NewEntryTX) (entrybus.NewEntry, error) {
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
	Data     *string `json:"data" validate:"required"`
	Metadata string  `json:"metadata" validate:"required"`
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

func toBusUpdateEntry(ctx context.Context, app UpdateEntry) (entrybus.UpdateEntry, error) {
	var e *entry.Entry
	if app.Data != nil {
		k, err := entry.Parse(*app.Data)
		if err != nil {
			return entrybus.UpdateEntry{}, fmt.Errorf("parse: %w", err)
		}
		e = &k
	}

	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return entrybus.UpdateEntry{}, err
	}

	bus := entrybus.UpdateEntry{
		Data:   e,
		UserID: &userID,
	}

	return bus, nil
}

// =============================================================================

// DeleteEntry defines the data needed to update a bundle metadata after deleting a password entry.
type DeleteEntry struct {
	Metadata string `json:"metadata" validate:"required"`
}

// Decode implements the decoder interface.
func (app *DeleteEntry) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Validate checks the data in the model is considered clean.
func (app DeleteEntry) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// =============================================================================
// Bundle

// Bundle represents information about an individual bundle.
type Bundle struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Type        string `json:"type"`
	Metadata    string `json:"metadata"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toBusUpdateBundle(metadata string) (bundlebus.UpdateBundle, error) {
	bus := bundlebus.UpdateBundle{
		Metadata: &metadata,
	}

	return bus, nil
}
