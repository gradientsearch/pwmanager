package vbundlebus_test

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/business/sdk/unitest"
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

	prds, err := keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu1 := unitest.User{
		User: usrs[0],
		Keys: prds,
	}

	// -------------------------------------------------------------------------

	usrs, err = userbus.TestSeedUsers(ctx, 1, role.Admin, busDomain.User)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding users : %w", err)
	}

	prds, err = keybus.TestGenerateSeedKeys(ctx, 2, busDomain.Key, usrs[0].ID)
	if err != nil {
		return unitest.SeedData{}, fmt.Errorf("seeding keys : %w", err)
	}

	tu2 := unitest.User{
		User: usrs[0],
		Keys: prds,
	}

	// -------------------------------------------------------------------------

	sd := unitest.SeedData{
		Admins: []unitest.User{tu2},
		Users:  []unitest.User{tu1},
	}

	return sd, nil
}

// =============================================================================

func toVBundle(usr userbus.User, prd keybus.Key) vbundlebus.Key {
	return vbundlebus.Key{
		ID:          prd.ID,
		UserID:      prd.UserID,
		Data:        prd.Data,
		DateCreated: prd.DateCreated,
		DateUpdated: prd.DateUpdated,
	}
}

func toVBundles(usr userbus.User, prds []keybus.Key) []vbundlebus.Key {
	items := make([]vbundlebus.Key, len(prds))
	for i, prd := range prds {
		items[i] = toVBundle(usr, prd)
	}

	return items
}

// =============================================================================

func query(busDomain dbtest.BusDomain, sd unitest.SeedData) []unitest.Table {
	prds := toVBundles(sd.Admins[0].User, sd.Admins[0].Keys)
	prds = append(prds, toVBundles(sd.Users[0].User, sd.Users[0].Keys)...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID.String() <= prds[j].ID.String()
	})

	table := []unitest.Table{
		{
			Name:    "all",
			ExpResp: prds,
			ExcFunc: func(ctx context.Context) any {
				filter := vbundlebus.QueryFilter{
					Name: dbtest.NamePointer("Name"),
				}

				resp, err := busDomain.VBundle.Query(ctx, filter, vbundlebus.DefaultOrderBy, page.MustParse("1", "10"))
				if err != nil {
					return err
				}

				return resp
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.([]vbundlebus.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.([]vbundlebus.Key)

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
