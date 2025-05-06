package vbundledb

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb/dbarray"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
)

func toBusBundleUsers(db string) ([]vbundlebus.BundleUser, error) {
	var users []bundleUser

	err := json.Unmarshal([]byte(db), &users)
	if err != nil {
		return nil, fmt.Errorf("unmarshal db bundle users")
	}

	bus := make([]vbundlebus.BundleUser, len(users))

	for i, u := range users {
		roles, err := bundlerole.ParseMany(u.Roles)
		if err != nil {
			return nil, fmt.Errorf("invalid users bundlerole :%s", err)
		}
		bus[i].Email = u.Email
		bus[i].Name = u.Name
		bus[i].UserID = u.UserID
		bus[i].Roles = roles
	}

	return bus, nil
}

func toBusBundle(db userBundleKey) (vbundlebus.UserBundleKey, error) {
	roles, err := bundlerole.ParseMany(db.KeyRoles)
	if err != nil {
		return vbundlebus.UserBundleKey{}, fmt.Errorf("invalid user bundlerole :%s", err)
	}

	users, err := toBusBundleUsers(db.Users)
	if err != nil {
		return vbundlebus.UserBundleKey{}, err
	}

	bus := vbundlebus.UserBundleKey{
		UserID:      db.UserID,
		Name:        db.Name,
		BundleID:    db.BundleID,
		Type:        db.Type,
		Metadata:    db.Metadata,
		KeyData:     db.KeyData,
		KeyRoles:    roles,
		Users:       users,
		DateCreated: db.DateCreated.In(time.Local),
		DateUpdated: db.DateUpdated.In(time.Local),
	}

	return bus, nil
}

func toBusBundles(dbKeys []userBundleKey) ([]vbundlebus.UserBundleKey, error) {
	bus := make([]vbundlebus.UserBundleKey, len(dbKeys))

	for i, dbKey := range dbKeys {
		var err error
		bus[i], err = toBusBundle(dbKey)
		if err != nil {
			return nil, err
		}
	}

	return bus, nil
}

// Nested user info inside the "users" JSON array
type bundleUser struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
	Roles  []string  `json:"roles"`
}

// Main structure for the query result
type userBundleKey struct {
	UserID      uuid.UUID      `db:"user_id"`
	Name        string         `db:"name"`
	BundleID    uuid.UUID      `db:"bundle_id"`
	Type        string         `db:"type"`
	Metadata    string         `db:"metadata"`
	DateCreated time.Time      `db:"date_created"`
	DateUpdated time.Time      `db:"date_updated"`
	KeyData     string         `db:"key_data"`
	KeyRoles    dbarray.String `db:"key_roles"`
	Users       string         `db:"users"`
}
