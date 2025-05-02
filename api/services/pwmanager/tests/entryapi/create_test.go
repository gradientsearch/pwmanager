package entry_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func create200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/bundles/" + sd.Users[0].Bundles[0].ID.String() + "/entries",
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &entryapp.NewEntryTX{
				Data:     "Guitar",
				Metadata: "UPDATED BUNDLE METADATA",
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					Data:     "Guitar",
					UserID:   sd.Users[0].ID.String(),
					BundleID: sd.Users[0].Bundles[0].ID.String(),
				},
				Bundle: entryapp.Bundle{
					Metadata: "UPDATED BUNDLE METADATA",
					UserID:   sd.Users[0].ID.String(),
					ID:       sd.Users[0].Bundles[0].ID.String(),
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
		},
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "missing-input",
			URL:        "/v1/bundles/" + sd.Users[0].Bundles[0].ID.String() + "/entries",
			Token:      sd.Users[0].Token,
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
			URL:        "/v1/bundles/" + sd.Users[0].Bundles[0].ID.String() + "/entries",
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
			URL:        "/v1/bundles/" + sd.Users[0].Bundles[0].ID.String() + "/entries",
			Token:      sd.Admins[0].Token[:10],
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
			URL:        "/v1/bundles/" + sd.Users[0].Bundles[0].ID.String() + "/entries",
			Token:      sd.Admins[0].Token + "A",
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        "/v1/bundles/" + sd.Admins[0].Bundles[0].ID.String() + "/entries",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[ADMIN]] rule[rule_user_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
