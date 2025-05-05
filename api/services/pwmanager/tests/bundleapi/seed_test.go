package bundle_test

import (
	"context"
	"fmt"

	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
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

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 5, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[userBundleAdmin].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	tu1 := apitest.User{
		User:    usrs[userBundleAdmin],
		Bundles: bdls,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userBundleAdmin].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu2

	tu2 := apitest.User{
		User:  usrs[userReadWrite],
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userReadWrite].Email.Address),
	}

	// -------------------------------------------------------------------------
	// tu3

	tu3 := apitest.User{
		User:    usrs[userRead],
		Bundles: bdls,
		Token:   apitest.Token(db.BusDomain.User, ath, usrs[userRead].Email.Address),
	}

	// -------------------------------------------------------------------------

	tu4 := apitest.User{
		User:  usrs[userNoRoles],
		Token: apitest.Token(db.BusDomain.User, ath, usrs[userNoRoles].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Users: []apitest.User{tu1, tu2, tu3, tu4},
	}

	return sd, nil
}
