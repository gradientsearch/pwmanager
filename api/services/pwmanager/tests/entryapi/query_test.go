package entry_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bundle-admin",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userBundleAdmin].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &entryapp.Entry{},
			ExpResp:    toAppEntryPtr(sd.Users[userBundleAdmin].Entries[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "shared-read-write",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userReadWrite].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &entryapp.Entry{},
			ExpResp:    toAppEntryPtr(sd.Users[userBundleAdmin].Entries[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "shared-user-read-only",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userRead].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &entryapp.Entry{},
			ExpResp:    toAppEntryPtr(sd.Users[userBundleAdmin].Entries[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}
	return table
}

func queryByID403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{

		{
			Name:       "shared-user-no-roles",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoRoles].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have read perms for bundle[%s] to read entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "shared-user-no-key",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoKey].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
	}
	return table
}
