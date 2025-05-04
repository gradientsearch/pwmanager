package entry_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func delete200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:  fmt.Sprintf("tu%d-user-bundle-admin", userBundleAdmin),
			URL:   fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token: sd.Users[userBundleAdmin].Token,
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			Method:     http.MethodDelete,
			StatusCode: http.StatusOK,
			GotResp:    &entryapp.EntryTx{},
			ExpResp:    &entryapp.EntryTx{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*entryapp.EntryTx)
				gotMetadata := gotResp.Bundle.Metadata
				expMetadata := "UPDATED BUNDLE METADATA"
				return cmp.Diff(gotMetadata, expMetadata)
			},
		},
		{
			Name:  fmt.Sprintf("tu%d-shared-read-write", userReadWrite),
			URL:   fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[1].ID),
			Token: sd.Users[userReadWrite].Token,
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			Method:     http.MethodDelete,
			StatusCode: http.StatusOK,
			GotResp:    &entryapp.EntryTx{},
			ExpResp:    &entryapp.EntryTx{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*entryapp.EntryTx)
				gotMetadata := gotResp.Bundle.Metadata
				expMetadata := "UPDATED BUNDLE METADATA"
				return cmp.Diff(gotMetadata, expMetadata)
			},
		},
	}

	return table
}

func delete401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[1].ID, sd.Users[userBundleAdmin].Entries[2].ID),
			Token:      "&nbsp;",
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[1].ID, sd.Users[userBundleAdmin].Entries[2].ID),
			Token:      sd.Users[userBundleAdmin].Token + "A",
			Method:     http.MethodDelete,
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

func delete403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-shared-user-read-only", userRead),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[2].ID),
			Token:      sd.Users[userRead].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-shared-user-no-roles", userNoRoles),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[2].ID),
			Token:      sd.Users[userNoRoles].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-shared-user-no-key", userNoKey),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[2].ID),
			Token:      sd.Users[userNoKey].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
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
