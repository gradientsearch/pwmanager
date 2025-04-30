package bundle_test

import (
	"fmt"
	"net/http"
	"time"

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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("SINGLE FAMILY"),
				Address: &bundleapp.UpdateAddress{
					Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
					Address2: dbtest.StringPointer("apt 105"),
					ZipCode:  dbtest.StringPointer("35810"),
					City:     dbtest.StringPointer("Huntsville"),
					State:    dbtest.StringPointer("AL"),
					Country:  dbtest.StringPointer("US"),
				},
			},
			GotResp: &bundleapp.Bundle{},
			ExpResp: &bundleapp.Bundle{
				ID:     sd.Users[0].Bundles[0].ID.String(),
				UserID: sd.Users[0].ID.String(),
				Type:   "SINGLE FAMILY",
				Address: bundleapp.Address{
					Address1: "123 Mocking Bird Lane",
					Address2: "apt 105",
					ZipCode:  "35810",
					City:     "Huntsville",
					State:    "AL",
					Country:  "US",
				},
				DateCreated: sd.Users[0].Bundles[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[0].Bundles[0].DateCreated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundleapp.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundleapp.Bundle)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-input",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.UpdateBundle{
				Address: &bundleapp.UpdateAddress{
					Address1: dbtest.StringPointer(""),
					Address2: dbtest.StringPointer(""),
					ZipCode:  dbtest.StringPointer(""),
					City:     dbtest.StringPointer(""),
					State:    dbtest.StringPointer(""),
					Country:  dbtest.StringPointer(""),
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"address1\",\"error\":\"address1 must be at least 1 character in length\"},{\"field\":\"zipCode\",\"error\":\"zipCode must be a valid numeric value\"},{\"field\":\"state\",\"error\":\"state must be at least 1 character in length\"},{\"field\":\"country\",\"error\":\"Key: 'UpdateBundle.address.country' Error:Field validation for 'country' failed on the 'iso3166_1_alpha2' tag\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-type",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.UpdateBundle{
				Type:    dbtest.StringPointer("BAD TYPE"),
				Address: &bundleapp.UpdateAddress{},
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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
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
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Admins[0].Bundles[0].ID),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("SINGLE FAMILY"),
				Address: &bundleapp.UpdateAddress{
					Address1: dbtest.StringPointer("123 Mocking Bird Lane"),
					Address2: dbtest.StringPointer("apt 105"),
					ZipCode:  dbtest.StringPointer("35810"),
					City:     dbtest.StringPointer("Huntsville"),
					State:    dbtest.StringPointer("AL"),
					Country:  dbtest.StringPointer("US"),
				},
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
