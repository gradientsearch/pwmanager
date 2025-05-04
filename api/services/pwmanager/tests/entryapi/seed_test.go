package entry_test

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
	userBundleAdmin userKey = iota
	userReadWrite
	userRead
	userNoRoles
	userNoKey
)

var userKeyMapping = map[userKey]string{
	userBundleAdmin: "user-bundle-owner",
	userReadWrite:   "user-read-write",
	userRead:        "user-read",
	userNoRoles:     "user-no-roles",
	userNoKey:       "user-no-key",
}

const (
	NUMBER_OF_USERS    = 5
	NUMBER_OF_BUNDLES  = 3
	ENTRIES_PER_BUNDLE = 10
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, NUMBER_OF_USERS, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// -------------------------------------------------------------------------
	// tu1

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, NUMBER_OF_BUNDLES, busDomain.Bundle, usrs[userBundleAdmin].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids := []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles := []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err := keybus.TestGenerateSeedKeys(ctx, NUMBER_OF_BUNDLES, busDomain.Key, usrs[userBundleAdmin].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	entries, err := entrybus.TestGenerateSeedEntries(ctx, ENTRIES_PER_BUNDLE, busDomain.Entry, usrs[userBundleAdmin].ID, bids)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding entries : %w", err)
	}

	tu1 := apitest.User{
		User:    usrs[userBundleAdmin],
		Keys:    keys,
		Bundles: bdls,
		Entries: entries,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userBundleAdmin].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu2

	roles = []bundlerole.Role{bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, NUMBER_OF_BUNDLES, busDomain.Key, usrs[userReadWrite].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := apitest.User{
		User:  usrs[userReadWrite],
		Keys:  keys,
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userReadWrite].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu3

	roles = []bundlerole.Role{bundlerole.Read}
	keys, err = keybus.TestGenerateSeedKeys(ctx, NUMBER_OF_BUNDLES, busDomain.Key, usrs[userRead].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu3 := apitest.User{
		User:  usrs[userRead],
		Keys:  keys,
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userRead].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu4

	roles = []bundlerole.Role{}
	keys, err = keybus.TestGenerateSeedKeys(ctx, NUMBER_OF_BUNDLES, busDomain.Key, usrs[userNoRoles].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu4 := apitest.User{
		User:  usrs[userNoRoles],
		Keys:  keys,
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userNoRoles].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu5

	tu5 := apitest.User{
		User:  usrs[userNoKey],
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userNoKey].Email.Address),
	}

	// -------------------------------------------------------------------------
	// admin

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// Create 3 bundles. Use last created bundle for create foreign key constraint
	// in create-200-basic test
	bdls, err = bundlebus.TestGenerateSeedBundles(ctx, 3, busDomain.Bundle, usrs[userBundleAdmin].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	bids = []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles = []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, NUMBER_OF_BUNDLES, busDomain.Key, usrs[userBundleAdmin].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	ta1 := apitest.User{
		User:    usrs[userBundleAdmin],
		Keys:    keys,
		Bundles: bdls,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userBundleAdmin].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{ta1},
		Users:  []apitest.User{tu1, tu2, tu3, tu4, tu5},
	}

	return sd, nil
}
