package keybus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/unitest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/key"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func Test_Key(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Key")

	sd, err := insertSeedData(db.BusDomain)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	unitest.Run(t, query(db.BusDomain, sd), "query")
	unitest.Run(t, create(db.BusDomain, sd), "create")
	unitest.Run(t, update(db.BusDomain, sd), "update")
	unitest.Run(t, delete(db.BusDomain, sd), "delete")
}

// =============================================================================

func insertSeedData(busDomain dbtest.BusDomain) (unitest.SeedData, error) {
	ctx := context.Background()

	usrs, err := userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err := bundlebus.TestGenerateSeedBundles(ctx, 3, busDomain.Bundle, usrs[0].ID)
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
		Keys:    keys,
		Bundles: bdls,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	bdls, err = bundlebus.TestGenerateSeedBundles(ctx, 3, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	bids = []uuid.UUID{}
	for _, v := range bdls {
		bids = append(bids, v.ID)
	}

	roles = []bundlerole.Role{bundlerole.Read, bundlerole.Write}
	keys, err = keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID, bids, roles)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := unitest.User{
		User:    usrs[0],
		Keys:    keys,
		Bundles: bdls,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Admins: []unitest.User{tu2},
		Users:  []unitest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	keys := make([]keybus.Key, 0, len(sd.Admins[0].Keys)+len(sd.Users[0].Keys))
	//keys = append(keys, sd.Admins[0].Keys...)
	keys = append(keys, sd.Users[0].Keys...)

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].ID.String() <= keys[j].ID.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: keys,
			ExcFunc: func(ctx context.Context) any {
				filter := keybus.QueryFilter{
					UserID: &sd.Users[0].ID,
				}

				resp, err := busDomain.Key.Query(ctx, filter, keybus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]keybus.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]keybus.Key)

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
		{
			Name:    "byid",
			ExpResp: sd.Users[0].Keys[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Key.QueryByID(ctx, sd.Users[0].Keys[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(keybus.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(keybus.Key)

				if gotResp.DateCreated.Format(time.RFC3339) == expResp.DateCreated.Format(time.RFC3339) {
					expResp.DateCreated = gotResp.DateCreated
				}

				if gotResp.DateUpdated.Format(time.RFC3339) == expResp.DateUpdated.Format(time.RFC3339) {
					expResp.DateUpdated = gotResp.DateUpdated
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func create(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: keybus.Key{
				UserID:   sd.Users[0].ID,
				BundleID: sd.Users[0].Bundles[2].ID,
				Data:     key.MustParse("Guitar"),
			},
			ExcFunc: func(ctx context.Context) any {
				nk := keybus.NewKey{
					UserID:   sd.Users[0].ID,
					BundleID: sd.Users[0].Bundles[2].ID,
					Data:     key.MustParse("Guitar"),
				}

				resp, err := busDomain.Key.Create(ctx, nk)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(keybus.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(keybus.Key)

				expResp.ID = gotResp.ID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name: "basic",
			ExpResp: keybus.Key{
				ID:          sd.Users[0].Keys[0].ID,
				BundleID:    sd.Users[0].Bundles[0].ID,
				UserID:      sd.Users[0].ID,
				Data:        key.MustParse("Guitar"),
				Roles:       []bundlerole.Role{bundlerole.Admin, bundlerole.Read, bundlerole.Write},
				DateCreated: sd.Users[0].Keys[0].DateCreated,
				DateUpdated: sd.Users[0].Keys[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uk := keybus.UpdateKey{
					Data: dbtest.KeyPointer("Guitar"),
				}

				resp, err := busDomain.Key.Update(ctx, sd.Users[0].Keys[0], uk)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(keybus.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(keybus.Key)

				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func delete(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	table := []unitest.Table{
		{
			Name:    "user",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Key.Delete(ctx, sd.Users[0].Keys[1]); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:    "admin",
			ExpResp: nil,
			ExcFunc: func(ctx context.Context) any {
				if err := busDomain.Key.Delete(ctx, sd.Admins[0].Keys[1]); err != nil {
					return err
				}

				return nil
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
