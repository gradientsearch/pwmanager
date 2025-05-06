package vbundlebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
)

// Nested user info inside the "users" JSON array
type BundleUser struct {
	UserID uuid.UUID
	Name   string
	Email  string
	Roles  []bundlerole.Role
}

// Main structure for the query result
type UserBundleKey struct {
	UserID      uuid.UUID
	Name        string
	BundleID    uuid.UUID
	Type        string
	Metadata    string
	DateCreated time.Time
	DateUpdated time.Time
	KeyData     string
	KeyRoles    []bundlerole.Role
	Users       []BundleUser
}
