package vbundlebus_test

import (
	"context"
	"fmt"
	"testing"
	"time"

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

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
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
		User: usrs[0],
		Keys: keys,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err = bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids = []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles = []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids, roles)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := unitest.User{
		User: usrs[0],
		Keys: keys,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Admins: []unitest.User{tu2},
		Users:  []unitest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func toVBundle(usr userbus.User, k keybus.Key) vbundlebus.UserBundleKey {
	return vbundlebus.UserBundleKey{
		UserID:      k.UserID,
		DateCreated: k.DateCreated,
		DateUpdated: k.DateUpdated,
	}
}

func toVBundles(usr userbus.User, keys []keybus.Key) []vbundlebus.UserBundleKey {
	items := make([]vbundlebus.UserBundleKey, len(keys))
	for i, k := range keys {
		items[i] = toVBundle(usr, k)
	}

	return items
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: []vbundlebus.UserBundleKey{},
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

				expResp := exp.([]vbundlebus.UserBundleKey)

				for i := range gotResp {
					if gotResp[i].DateCreated.Format(time.RFC3339) == expResp[i].DateCreated.Format(time.RFC3339) {
						expResp[i].DateCreated = gotResp[i].DateCreated
					}

					if gotResp[i].DateUpdated.Format(time.RFC3339) == expResp[i].DateUpdated.Format(time.RFC3339) {
						expResp[i].DateUpdated = gotResp[i].DateUpdated
					}
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
