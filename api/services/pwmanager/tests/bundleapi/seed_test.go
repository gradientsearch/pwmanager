package bundle_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

type userKey int

const (
	userA userKey = iota
	userB
)

var userKeyMapping = map[userKey]string{
	userA: "user-a",
	userB: "user-b",
}

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[userA].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	tu1 := apitest.User{
		User:    usrs[userA],
		Bundles: bdls,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userA].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu2 check cascading deletes

	bdls, err = bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[userB].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids := []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles := []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[userB].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	entries, err := entrybus.TestGenerateSeedEntries(ctx, 1, busDomain.Entry, usrs[userB].ID, bids)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding entries : %w", err)
	}

	tu2 := apitest.User{
		User:    usrs[userB],
		Bundles: bdls,
		Entries: entries,
		Keys:    keys,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userB].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Users: []apitest.User{tu1, tu2},
	}

	return sd, nil
}
