package entry_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
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
			Name:       "userrw",
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
			Name:       "shared-readonly",
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
		{
			Name:       "shared-noroles",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoRoles].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &entryapp.Entry{},
			ExpResp:    toAppEntryPtr(sd.Users[userBundleAdmin].Entries[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "shared-nokey",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[userNoRoles].Token,
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
