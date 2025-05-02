// Package mid provides app level middleware support.
package mid

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// isError tests if the Encoder has an error inside of it.
func isError(e web.Encoder) error {
	err, isError := e.(error)
	if isError {
		return err
	}
	return nil
}

// =============================================================================

type ctxKey int

const (
	claimKey ctxKey = iota + 1
	userIDKey
	userKey
	keyKey
	entryKey
	bundleKey
	trKey
)

func setClaims(ctx context.Context, claims auth.Claims) context.Context {
	return context.WithValue(ctx, claimKey, claims)
}

// GetClaims returns the claims from the context.
func GetClaims(ctx context.Context) auth.Claims {
	v, ok := ctx.Value(claimKey).(auth.Claims)
	if !ok {
		return auth.Claims{}
	}
	return v
}

func setUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserID returns the user id from the context.
func GetUserID(ctx context.Context) (uuid.UUID, error) {
	v, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("user id not found in context")
	}

	return v, nil
}

func setUser(ctx context.Context, usr userbus.User) context.Context {
	return context.WithValue(ctx, userKey, usr)
}

// GetUser returns the user from the context.
func GetUser(ctx context.Context) (userbus.User, error) {
	v, ok := ctx.Value(userKey).(userbus.User)
	if !ok {
		return userbus.User{}, errors.New("user not found in context")
	}

	return v, nil
}

func setKey(ctx context.Context, k keybus.Key) context.Context {
	return context.WithValue(ctx, keyKey, k)
}

// GetKey returns the key from the context.
func GetKey(ctx context.Context) (keybus.Key, error) {
	v, ok := ctx.Value(keyKey).(keybus.Key)
	if !ok {
		return keybus.Key{}, errors.New("key not found in context")
	}

	return v, nil
}

func setEntry(ctx context.Context, k entrybus.Entry) context.Context {
	return context.WithValue(ctx, entryKey, k)
}

// GetEntry returns the key from the context.
func GetEntry(ctx context.Context) (entrybus.Entry, error) {
	v, ok := ctx.Value(entryKey).(entrybus.Entry)
	if !ok {
		return entrybus.Entry{}, errors.New("key not found in context")
	}

	return v, nil
}

func setBundle(ctx context.Context, bdl bundlebus.Bundle) context.Context {
	return context.WithValue(ctx, bundleKey, bdl)
}

// GetBundle returns the bundle from the context.
func GetBundle(ctx context.Context) (bundlebus.Bundle, error) {
	v, ok := ctx.Value(bundleKey).(bundlebus.Bundle)
	if !ok {
		return bundlebus.Bundle{}, errors.New("bundle not found in context")
	}

	return v, nil
}

func setTran(ctx context.Context, tx sqldb.CommitRollbacker) context.Context {
	return context.WithValue(ctx, trKey, tx)
}

// GetTran retrieves the value that can manage a transaction.
func GetTran(ctx context.Context) (sqldb.CommitRollbacker, error) {
	v, ok := ctx.Value(trKey).(sqldb.CommitRollbacker)
	if !ok {
		return nil, errors.New("transaction not found in context")
	}

	return v, nil
}
