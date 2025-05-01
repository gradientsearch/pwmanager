package bundletxapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/types/key"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

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

// NewTran represents an example of cross domain transaction at the
// application layer.
type NewTran struct {
	Key  NewKey  `json:"key"`
	User NewUser `json:"user"`
}

// Validate checks the data in the model is considered clean.
func (app NewTran) Validate() error {
	if err := errs.Check(app); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

// Decode implements the decoder interface.
func (app *NewTran) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

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

func toBusNewKey(app NewKey) (keybus.NewKey, error) {
	k, err := key.Parse(app.Data)
	if err != nil {
		return keybus.NewKey{}, fmt.Errorf("parse: %w", err)
	}

	bus := keybus.NewKey{
		Data: k,
	}

	return bus, nil
}
