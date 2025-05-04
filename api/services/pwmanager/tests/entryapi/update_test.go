package entry_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-user-bundle-admin", userBundleAdmin),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar2"),
				Metadata: *dbtest.StringPointer("UPDATED METADATA"),
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					ID:          sd.Users[userBundleAdmin].Entries[0].ID.String(),
					UserID:      sd.Users[userBundleAdmin].ID.String(),
					BundleID:    sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					Data:        "Guitar2",
					DateCreated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
				},
				Bundle: entryapp.Bundle{
					ID:       sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					UserID:   sd.Users[userBundleAdmin].Bundles[0].UserID.String(),
					Type:     sd.Users[userBundleAdmin].Bundles[0].Type.String(),
					Metadata: "UPDATED METADATA",

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
		},
		{
			Name:       fmt.Sprintf("tu%d-shared-shared-user-read-write", userReadWrite),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userReadWrite].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar2"),
				Metadata: *dbtest.StringPointer("UPDATED METADATA"),
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					ID:          sd.Users[userBundleAdmin].Entries[0].ID.String(),
					UserID:      sd.Users[userReadWrite].ID.String(),
					BundleID:    sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					Data:        "Guitar2",
					DateCreated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[userBundleAdmin].Entries[0].DateCreated.Format(time.RFC3339),
				},
				Bundle: entryapp.Bundle{
					ID:       sd.Users[userBundleAdmin].Bundles[0].ID.String(),
					UserID:   sd.Users[userBundleAdmin].Bundles[0].UserID.String(),
					Type:     sd.Users[userBundleAdmin].Bundles[0].Type.String(),
					Metadata: "UPDATED METADATA",

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
		},
	}

	return table
}

func update401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      "&nbsp;",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userBundleAdmin].Token + "A",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-shared-user-read-only", userRead),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Entries[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userRead].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar"),
				Metadata: "New METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-shared-user-no-roles", userNoRoles),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Entries[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoRoles].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar"),
				Metadata: "New METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-user-no-key", userNoKey),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Entries[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoKey].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar"),
				Metadata: "New METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
