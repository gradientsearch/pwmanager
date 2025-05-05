package bundleapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
	"github.com/gradientsearch/pwmanager/business/types/name"

	"github.com/gradientsearch/pwmanager/business/types/key"

	"github.com/gradientsearch/pwmanager/business/types/role"
)

// Encode implements the encoder interface.
func (app Bundle) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppBundles(bundles []bundlebus.Bundle) []Bundle {
	app := make([]Bundle, len(bundles))
	for i, b := range bundles {
		app[i] = toAppBundle(b)
	}

	return app
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

// Key represents an individual key.
type Key struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Data        string `json:"data"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
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
		Data:        k.Data.String(),
		DateCreated: k.DateCreated.Format(time.RFC3339),
		DateUpdated: k.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

// NewBundleTx represents an example of cross domain transaction at the
// application layer.
type NewBundleTx struct {
	Bundle NewBundle `json:"bundle" validate:"required"`
	Key    NewKey    `json:"key" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (app NewBundleTx) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Decode implements the decoder interface.
func (app *NewBundleTx) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// =============================================================================

// BundleTx defines the data needed associated with a bundle transaction.
type BundleTx struct {
	Bundle Bundle
	Key    Key
}

func toAppBundleTx(b bundlebus.Bundle, k keybus.Key) BundleTx {
	return BundleTx{
		Bundle: toAppBundle(b),
		Key:    toAppKey(k),
	}
}

// Encode implements the encoder interface.
func (app BundleTx) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

// =============================================================================

// =============================================================================
// NewUser contains information needed to create a new user.
type NewUser struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Roles           []string `json:"roles" validate:"required"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"eqfield=Password"`
}

// Validate checks the data in the model is considered clean.
func (app NewUser) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewUser(app NewUser) (userbus.NewUser, error) {
	roles, err := role.ParseMany(app.Roles)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(app.Email)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(app.Name)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	department, err := name.ParseNull(app.Department)
	if err != nil {
		return userbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := userbus.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   app.Password,
	}

	return bus, nil
}

// =============================================================================

// NewKey is what we require from clients when adding a Key.
type NewKey struct {
	Data string `json:"data" validate:"required"`
}

// Validate checks the data in the model is considered clean.
func (app NewKey) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusNewKey(ctx context.Context, app NewKey, bundleID uuid.UUID) (keybus.NewKey, error) {
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("getuserid: %w", err)
	}

	k, err := key.Parse(app.Data)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse: %w", err)
	}

	bus := keybus.NewKey{
		UserID:   userID,
		BundleID: bundleID,
		Data:     k,
	}

	return bus, nil
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
// Bundle defines the data needed to add a new bundle.
type Bundle struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Type        string `json:"type"`
	Metadata    string `json:"metadata"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

func toAppBundle(b bundlebus.Bundle) Bundle {
	return Bundle{
		ID:          b.ID.String(),
		UserID:      b.UserID.String(),
		Type:        b.Type.String(),
		Metadata:    b.Metadata,
		DateCreated: b.DateCreated.String(),
		DateUpdated: b.DateUpdated.String(),
	}
}
