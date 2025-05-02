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
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar2"),
				Metadata: *dbtest.StringPointer("UPDATED METADATA"),
			},
			GotResp: &entryapp.EntryTx{},
			ExpResp: &entryapp.EntryTx{
				Entry: entryapp.Entry{
					ID:          sd.Users[0].Entries[0].ID.String(),
					UserID:      sd.Users[0].ID.String(),
					BundleID:    sd.Users[0].Bundles[0].ID.String(),
					Data:        "Guitar2",
					DateCreated: sd.Users[0].Entries[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[0].Entries[0].DateCreated.Format(time.RFC3339),
				},
				Bundle: entryapp.Bundle{
					ID:       sd.Users[0].Bundles[0].ID.String(),
					UserID:   sd.Users[0].Bundles[0].UserID.String(),
					Type:     sd.Users[0].Bundles[0].Type.String(),
					Metadata: "UPDATED METADATA",

					DateCreated: sd.Users[0].Bundles[0].DateCreated.Format(time.RFC3339),
					DateUpdated: sd.Users[0].Bundles[0].DateUpdated.Format(time.RFC3339),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[0].ID),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[0].ID),
			Token:      sd.Users[0].Token + "A",
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Entries[0].ID, sd.Users[0].Entries[0].ID),
			Token:      sd.Users[1].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar"),
				Metadata: "New METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: &errs.Error{},
			CmpFunc: func(got any, exp any) string {
				return ""
			},
		},
	}

	return table
}

func update403(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Entries[0].ID, sd.Users[0].Entries[0].ID),
			Token:      sd.Users[1].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &entryapp.UpdateEntry{
				Data:     dbtest.StringPointer("Guitar"),
				Metadata: "New METADATA",
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return ""
			},
		},
	}

	return table
}
