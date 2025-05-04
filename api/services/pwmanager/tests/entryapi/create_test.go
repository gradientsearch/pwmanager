package entry_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

func create200(sd apitest.SeedData) []apitest.Table {
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &entryapp.NewEntryTX{
				Data:     fmt.Sprintf("DATA%d", i.user),
				Metadata: fmt.Sprintf("METADATA%d", i.user),
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					Data:     fmt.Sprintf("DATA%d", i.user),
					UserID:   sd.Users[i.user].ID.String(),
					BundleID: sd.Users[userBundleAdmin].Bundles[0].ID.String(),
				},
				Bundle: entryapp.Bundle{
					Metadata: fmt.Sprintf("METADATA%d", i.user),
					Type:     bundletype.Shareable.String(),
					UserID:   sd.Users[userBundleAdmin].ID.String(),
					ID:       sd.Users[userBundleAdmin].Bundles[0].ID.String(),
				},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*entryapp.EntryTx)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*entryapp.EntryTx)
				expResp.Entry.ID = gotResp.Entry.ID
				expResp.Entry.DateCreated = gotResp.Entry.DateCreated
				expResp.Entry.DateUpdated = gotResp.Entry.DateUpdated
				expResp.Bundle.Type = gotResp.Bundle.Type
				expResp.Bundle.DateCreated = gotResp.Bundle.DateCreated
				expResp.Bundle.DateUpdated = gotResp.Bundle.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		}
		table = append(table, t)
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &entryapp.NewEntryTX{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"data\",\"error\":\"data is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      "&nbsp;",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badtoken",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userBundleAdmin].Token[:10],
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "badsig",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userBundleAdmin].Token + "A",
			Method:     http.MethodPost,
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

func create403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-shared-read-only", userRead),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userRead].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.NewEntryTX{
				Data:     "Guitar",
				Metadata: "UPDATED BUNDLE METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				// &errs.Error{Code:errs.ErrCode{value:8}, Message:"must have write perms for bundle[bb459b0a-79a5-4af9-b89a-67d830ec1db9] to create an entry", FuncName:"", FileName:""}
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID)
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-shared-no-roles", userNoRoles),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userNoRoles].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.NewEntryTX{
				Data:     "Guitar",
				Metadata: "UPDATED BUNDLE METADATA",
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
			Name:       fmt.Sprintf("tu%d-shared-no-key", userNoKey),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[userNoKey].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.NewEntryTX{
				Data:     "Guitar",
				Metadata: "UPDATED BUNDLE METADATA",
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
