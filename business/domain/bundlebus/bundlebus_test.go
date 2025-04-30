package bundlebus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/unitest"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
	"github.com/gradientsearch/pwmanager/business/types/role"
)

func Test_Bundle(t *testing.T) {
	t.Parallel()

	db := dbtest.New(t, "Test_Bundle")

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

	hmes, err := bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	tu1 := unitest.User{
		User:    usrs[0],
		Bundles: hmes,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.User, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu2 := unitest.User{
		User: usrs[0],
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	hmes, err = bundlebus.TestGenerateSeedBundles(ctx, 2, busDomain.Bundle, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding bundles : %w", err)
	}

	tu3 := unitest.User{
		User:    usrs[0],
		Bundles: hmes,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	tu4 := unitest.User{
		User: usrs[0],
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Users:  []unitest.User{tu1, tu2},
		Admins: []unitest.User{tu3, tu4},
	}

	return sd, nil
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	hmes := make([]bundlebus.Bundle, 0, len(sd.Admins[0].Bundles)+len(sd.Users[0].Bundles))
	hmes = append(hmes, sd.Admins[0].Bundles...)
	hmes = append(hmes, sd.Users[0].Bundles...)

	sort.Slice(hmes, func(i, j int) bool {
		return hmes[i].ID.String() <= hmes[j].ID.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: hmes,
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Bundle.Query(ctx, bundlebus.QueryFilter{}, bundlebus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]bundlebus.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]bundlebus.Bundle)

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
			ExpResp: sd.Users[0].Bundles[0],
			ExcFunc: func(ctx context.Context) any {
				resp, err := busDomain.Bundle.QueryByID(ctx, sd.Users[0].Bundles[0].ID)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(bundlebus.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(bundlebus.Bundle)

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
			ExpResp: bundlebus.Bundle{
				UserID: sd.Users[0].ID,
				Type:   bundletype.Single,
				Address: bundlebus.Address{
					Address1: "123 Mocking Bird Lane",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
			},
			ExcFunc: func(ctx context.Context) any {
				nh := bundlebus.NewBundle{
					UserID: sd.Users[0].ID,
					Type:   bundletype.Single,
					Address: bundlebus.Address{
						Address1: "123 Mocking Bird Lane",
						ZipCode:  "35810",
						City:     "Huntsville",
						State:    "AL",
						Country:  "US",
					},
				}

				resp, err := busDomain.Bundle.Create(ctx, nh)
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(bundlebus.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(bundlebus.Bundle)

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
			ExpResp: bundlebus.Bundle{
				ID:     sd.Users[0].Bundles[0].ID,
				UserID: sd.Users[0].ID,
				Type:   bundletype.Single,
				Address: bundlebus.Address{
					Address1: "123 Mocking Bird Lane",
					Address2: "apt 105",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
				DateCreated: sd.Users[0].Bundles[0].DateCreated,
				DateUpdated: sd.Users[0].Bundles[0].DateCreated,
			},
			ExcFunc: func(ctx context.Context) any {
				uh := bundlebus.UpdateBundle{
					Type: &bundletype.Single,
					Address: &bundlebus.UpdateAddress{
						Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
						Address2: dbtest.StringPointer("apt 105"),
						ZipCode:  dbtest.StringPointer("35810"),
						City:     dbtest.StringPointer("Huntsville"),
						State:    dbtest.StringPointer("AL"),
						Country:  dbtest.StringPointer("US"),
					},
				}

				resp, err := busDomain.Bundle.Update(ctx, sd.Users[0].Bundles[0], uh)
				if err != nil {
					return err
				}

				resp.DateUpdated = resp.DateCreated

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(bundlebus.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(bundlebus.Bundle)

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
				if err := busDomain.Bundle.Delete(ctx, sd.Users[0].Bundles[1]); err != nil {
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
				if err := busDomain.Bundle.Delete(ctx, sd.Admins[0].Bundles[1]); err != nil {
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
