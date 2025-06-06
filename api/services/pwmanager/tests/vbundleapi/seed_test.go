package vbundle_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids := []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles := []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu1 := apitest.User{
		User:    usrs[0],
		Bundles: bdls,
		Keys:    keys,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	roles = []bundlerole.Role{bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, 1, busDomain.Key, usrs[1].ID, []uuid.UUID{bids[0]}, roles)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := apitest.User{
		User:  usrs[1],
		Keys:  keys,
		Token: apitest.Token(db.BusDomain.User, ath, usrs[1].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Users: []apitest.User{tu1, tu2},
	}

	return sd, nil
}
