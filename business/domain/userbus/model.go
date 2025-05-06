package userbus

import (
	"net/mail"
	"time"

	"github.com/google/uuid"
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
	UUK          UUK
}

// NewUser contains information needed to create a new user.
type NewUser struct {
	Name       name.Name
	Email      mail.Address
	Roles      []role.Role
	Department name.Null
	Password   string
	UUK        UUK
}

// UpdateUser contains information needed to update a user.
type UpdateUser struct {
	Name       *name.Name
	Email      *mail.Address
	Roles      []role.Role
	Department *name.Null
	Password   *string
	Enabled    *bool
	UUK        *UUK
}

// -------------------------------------------------------------------------
// Encryption structs
type EncSymKey struct {
	// uuid of the private key
	Kid string `json:"kid"`
	// encoding used to encrypt the data e.g. A256GCM
	Enc string `json:"enc"`
	// initialization
	Iv string `json:"iv"`
	// encrypted symmetric key
	Data string `json:"data"`
	// content type
	Cty string `json:"cty"`
	// the algorithm used to encrypt the EncSymKey e.g. 2SKD PBDKF2-HKDF
	Alg string `json:"alg"`
	// PBDKF2 iterations e.g. 650000
	P2c int `json:"p2c"`
	// initial 16 byte random sequence for secret key derivation.
	// used in the first hkdf function call
	P2s string `json:"p2s"`
}

type EncPriKey struct {
	// uuid
	Kid string `json:"kid"`
	// encoding of data e.g. A256GCM
	Enc string `json:"enc"`
	// initialization vector used to encrypt the priv key
	Iv string `json:"iv"`
	// the encrypted priv key
	Data string `json:"data"`
	// format used for encrypted data e.g JWK format
	Cty string `json:"cty"`
}

// user unlock key
// The secret key encrypts the EncSymKey, the EncSymKey
// encrypts the users PrivateKey
type UUK struct {
	// uuid of priv key
	UUID string `json:"uuid"`
	// symmetric key used to encrypt the EncPriKey
	EncSymKey EncSymKey `json:"enc_sym_key"`
	// mp a.k.a secret key
	EncryptedBy string `json:"encrypted_by"`
	// priv key used to encrypt `Safe` data
	EncPriKey EncPriKey `json:"enc_pri_key"`
	// pub key of the private key
	PubKey interface{} `json:"pub_key"`
}
