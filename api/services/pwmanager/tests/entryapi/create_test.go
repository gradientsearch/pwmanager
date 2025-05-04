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
			Name:       fmt.Sprintf("tu%d-missing-input", userBundleAdmin),
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
	table := []apitest.Table{}
	inputs := []struct {
		name  string
		token string
		err   *errs.Error
	}{
		{
			"emptytoken",
			"&nbsp;",
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			"badtoken",
			sd.Users[userBundleAdmin].Token[:10],
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			"badsig",
			sd.Users[userBundleAdmin].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", userBundleAdmin, i.name),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      i.token,
			Method:     http.MethodPost,
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

func create403(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user       userKey
		errMessage string
	}{
		{
			userRead,
			fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoRoles,
			fmt.Sprintf("must have write perms for bundle[%s] to create an entry", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoKey,
			fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", userRead, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[i.user].Token,
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
				expResp.Message = i.errMessage
				return cmp.Diff(got, exp)
			},
		}

		table = append(table, t)
	}
	return table
}
