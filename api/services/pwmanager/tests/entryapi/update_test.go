package entry_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func update200(sd apitest.SeedData) []apitest.Table {

	inputs := []struct {
		user userKey
	}{
		{
			userBundleAdmin,
		},
		{
			userReadWrite,
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &entryapp.UpdateEntry{
				Data:     fmt.Sprintf("%s%d", "Guitar", i.user),
				Metadata: fmt.Sprintf("%s%d", "Metadata", i.user),
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					ID:          sd.Users[userBundleAdmin].Entries[0].ID.String(),
					UserID:      sd.Users[i.user].ID.String(),
					BundleID:    sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					Data:        fmt.Sprintf("%s%d", "Guitar", i.user),
					DateCreated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
				},
				Bundle: entryapp.Bundle{
					ID:       sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					UserID:   sd.Users[userBundleAdmin].Bundles[0].UserID.String(),
					Type:     sd.Users[userBundleAdmin].Bundles[0].Type.String(),
					Metadata: fmt.Sprintf("%s%d", "Metadata", i.user),

					DateCreated: sd.Users[userBundleAdmin].Bundles[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[userBundleAdmin].Bundles[0].DateUpdated.Format(time.RFC3339),
				},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*entryapp.EntryTx)

				if !exists {
					return "error occurred"
				}

				expResp := exp.(*entryapp.EntryTx)
				gotResp.Bundle.DateUpdated = expResp.Bundle.DateUpdated
				gotResp.Entry.DateUpdated = expResp.Entry.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		}
		table = append(table, t)
	}
	return table
}

func update401(sd apitest.SeedData) []apitest.Table {

	inputs := []struct {
		user  userKey
		name  string
		token string
		err   *errs.Error
	}{
		{
			userBundleAdmin,
			"emptytoken",
			"&nbsp;",
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			userBundleAdmin,
			"badsig",
			sd.Users[userBundleAdmin].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, i.name),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      i.token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    i.err,
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		}
		table = append(table, t)
	}

	return table
}

func update403(sd apitest.SeedData) []apitest.Table {
	permError := fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)

	roles := []struct {
		user       userKey
		errMessage string
	}{

		{
			userRead,
			permError,
		},
		{
			userNoRoles,
			permError,
		},
		{
			userNoKey,
			fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID),
		},
	}

	table := []apitest.Table{}
	for _, r := range roles {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", r.user, userKeyMapping[r.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Entries[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[r.user].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     "Guitar",
				Metadata: "NEW METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = r.errMessage
				return cmp.Diff(got, exp)
			},
		}
		table = append(table, t)
	}

	return table
}
