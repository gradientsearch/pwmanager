package entry_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func delete200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:  "asuser",
			URL:   fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[0].ID),
			Token: sd.Users[0].Token,
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[1].ID, sd.Users[0].Entries[1].ID),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[1].ID, sd.Users[0].Entries[1].ID),
			Token:      sd.Users[0].Token + "A",
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "wronguser",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[1].ID),
			Token:      sd.Users[1].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			Input: &entryapp.DeleteEntry{
				Metadata: "UPDATED BUNDLE METADATA",
			},
			ExpResp: errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[ADMIN]] rule[rule_user_only]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*errs.Error)
				if strings.Contains(gotResp.Message, "db: entry not found") {
					return ""
				}
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "asadmin",
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[0].Bundles[0].ID, sd.Users[0].Entries[0].ID),
			Token:      sd.Admins[0].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    &errs.Error{},
			CmpFunc: func(got any, exp any) string {
				return ""
			},
		},
	}
	return table
}
