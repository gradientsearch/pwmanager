package vbundlebus_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/sdk/unitest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func Test_VBundle(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Key")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unitest.Run(t, query(db.BusDomain, sd), "query")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := userbus.TestSeedUsers(ctx, 2, role.User, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids := []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles := []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids, roles)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu1 := unitest.User{
		User:    usrs[0],
		Bundles: bdls,
		Keys:    keys,
	}

	// -------------------------------------------------------------------------

	roles = []bundlerole.Role{bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, 1, busDomain.Key, usrs[1].ID, []uuid.UUID{bids[0]}, roles)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := unitest.User{
		User: usrs[1],
		Keys: keys,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Users: []unitest.User{tu1, tu2},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: make([]vbundlebus.UserBundleKey, 2),
			ExcFunc: func(ctx context.Context) any {

				resp, err := busDomain.VBundle.QueryByID(ctx, sd.Users[0].User.ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]vbundlebus.UserBundleKey)
				if !exists {
					return "error occurred"
				}

				tu1 := sd.Users[0]
				tu2 := sd.Users[1]
				// =============================================================================

				b1 := gotResp[0]
				b2 := gotResp[1]

				expResp := []vbundlebus.UserBundleKey{
					{
						UserID:      tu1.User.ID,
						BundleID:    tu1.Bundles[0].ID,
						Name:        tu1.User.Name.String(),
						Type:        b1.Type,
						Metadata:    tu1.Bundles[0].Metadata,
						DateCreated: b1.DateCreated,
						DateUpdated: b1.DateUpdated,
						KeyData:     tu1.Keys[0].Data.String(),
						KeyRoles:    tu1.Keys[0].Roles,
						Users: []vbundlebus.BundleUser{
							{
								UserID: tu1.User.ID,
								Name:   tu1.User.Name.String(),
								Email:  tu1.User.Email.Address,
								Roles:  tu1.Keys[0].Roles,
							},
							{
								UserID: tu2.User.ID,
								Name:   tu2.User.Name.String(),
								Email:  tu2.User.Email.Address,
								Roles:  tu2.Keys[0].Roles,
							},
						},
					},
					{
						UserID:      tu1.User.ID,
						BundleID:    tu1.Bundles[1].ID,
						Name:        tu1.User.Name.String(),
						Type:        b2.Type,
						Metadata:    tu1.Bundles[1].Metadata,
						DateCreated: b2.DateCreated,
						DateUpdated: b2.DateUpdated,
						KeyData:     tu1.Keys[1].Data.String(),
						KeyRoles:    tu1.Keys[1].Roles,
						Users: []vbundlebus.BundleUser{
							{
								UserID: tu1.User.ID,
								Name:   tu1.User.Name.String(),
								Email:  tu1.User.Email.Address,
								Roles:  tu1.Keys[1].Roles,
							},
						},
					},
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
