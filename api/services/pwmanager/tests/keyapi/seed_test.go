package key_test

import (
	"context"
	"fmt"

	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func insertSeedData(db *dbtest.Database, ath *auth.Auth) (apitest.SeedData, error) {
	ctx := context.Background()
	busDomain := db.BusDomain

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu1 := apitest.User{
		User: usrs[0],
		Keys: prds,

		Token: apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID)
	if err != nil {
		return apitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := apitest.User{
		User:  usrs[0],
		Keys:  prds,
		Token: apitest.Token(db.BusDomain.User, ath, usrs[0].Email.Address),
	}

	// -------------------------------------------------------------------------

	sd := apitest.SeedData{
		Admins: []apitest.User{tu2},
		Users:  []apitest.User{tu1},
	}

	return sd, nil
}
