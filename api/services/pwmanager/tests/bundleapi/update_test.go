package bundle_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
			Token:      sd.Users[userA].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("PERSONAL"),
			},
			GotResp: &bundleapp.Bundle{},
			ExpResp: &bundleapp.Bundle{
				ID:     sd.Users[userA].Bundles[0].ID.String(),
				UserID: sd.Users[userA].ID.String(),
				Type:   "PERSONAL",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundleapp.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundleapp.Bundle)
				gotResp.DateUpdated = expResp.DateUpdated
				gotResp.DateCreated = expResp.DateCreated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-type",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
			Token:      sd.Users[userA].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("BAD TYPE"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "parse: invalid bundle type \"BAD TYPE\""),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func update401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "emptytoken",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
			Token:      sd.Users[userA].Token + "A",
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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
			Token:      sd.Users[userB].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("PERSONAL"),
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.Unauthenticated, "authorize: you are not authorized for that action, claims[[USER]] rule[rule_admin_or_subject]: rego evaluation failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
