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
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 3, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids := []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	keys, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	entries, err := entrybus.TestGenerateSeedEntries(ctx, 2, busDomain.Entry, usrs[0].ID, bids)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding entries : %w", err)
	}

	tu1 := apitest.User{
		User:    usrs[0],
		Keys:    keys,
		Bundles: bdls,
		Entries: entries,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	tuOther := apitest.User{
		User:  usrs[1],
		Token: apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	// Create 3 bundles. Use last created bundle for create foreign key constraint
	// in create-200-basic test
	bdls, err = bundlebus.TestGenerateSeedBundles(ctx, 3, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	bids = []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	keys, err = keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := apitest.User{
		User:    usrs[0],
		Keys:    keys,
		Bundles: bdls,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu2},
		Users:  []apitest.User{tu1, tuOther},
	}

	return sd, nil
}
