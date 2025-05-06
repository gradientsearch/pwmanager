package vbundleapp

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
)

// Nested user info inside the "users" JSON array
type BundleUser struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
}

// Main structure for the query result
type UserBundleKey struct {
	UserID      uuid.UUID    `json:"user_id"`
	Name        string       `json:"name"`
	BundleID    uuid.UUID    `json:"bundle_id"`
	Type        string       `json:"type"`
	Metadata    string       `json:"metadata"`
	DateCreated time.Time    `json:"date_created"`
	DateUpdated time.Time    `json:"date_updated"`
	KeyData     string       `json:"key_data"`
	KeyRoles    []string     `json:"key_roles"`
	Users       []BundleUser `json:"users"`
}

type UserBundleKeys []UserBundleKey

// Decode implements the decoder interface.
func (app *UserBundleKeys) Decode(data []byte) error {
	return json.Unmarshal(data, app)
}

// Encode implements the encoder interface.
func (app UserBundleKeys) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppUserBundleKey(ub vbundlebus.UserBundleKey) UserBundleKey {
	return UserBundleKey{
		UserID:      ub.UserID,
		KeyData:     ub.KeyData,
		DateCreated: ub.DateCreated,
		DateUpdated: ub.DateUpdated,
	}
}

func toAppUserBundleKeys(keys []vbundlebus.UserBundleKey) UserBundleKeys {
	app := make(UserBundleKeys, len(keys))
	for i, k := range keys {
		app[i] = toAppUserBundleKey(k)
	}

	return app
}
