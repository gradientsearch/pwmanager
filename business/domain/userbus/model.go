package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/sdk/uuk"
	"github.com/gradientsearch/pwmanager/business/types/name"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID           uuid.UUID
	Name         name.Name
	Email        mail.Address
	Roles        []role.Role
	PasswordHash []byte
	Department   name.Null
	Enabled      bool
	DateCreated  time.Time
	DateUpdated  time.Time
	UUK          uuk.UUK
}

// NewUser contains information needed to create a new user.
// New users are created by admins. Users will use Register
// to update the password and create a UUK payload.
type NewUser struct {
	Name       name.Name
	Email      mail.Address
	Roles      []role.Role
	Department name.Null
	Password   string
}

// RegisterUser contains information needed to register a user with a password and UUK.
// Token is the users register token provided by an admin.
type RegisterUser struct {
	Token    string
	Password string
	UUK      uuk.UUK
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name       *name.Name
	Email      *mail.Address
	Roles      []role.Role
	Department *name.Null
	Password   *string
	Enabled    *bool
}

// UpdateUserPassword contains information needed to update user password.
type UpdateUserPassword struct {
	Password *string
}
