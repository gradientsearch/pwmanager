package key_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func create200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[1].ID.String()),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &keyapp.NewKey{
				BundleID: sd.Users[0].Bundles[2].ID.String(),
				UserID:   string(sd.Users[1].ID[0]),
				Data:     "Guitar",
				Roles:    []string{"ADMIN", "WRITE", "READ"},
			},
			GotResp: &keyapp.Key{},
			ExpResp: &keyapp.Key{
				Data:     "Guitar",
				UserID:   sd.Users[0].ID.String(),
				BundleID: sd.Users[0].Bundles[2].ID.String(),
				Roles:    []string{"ADMIN", "WRITE", "READ"},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*keyapp.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*keyapp.Key)

				expResp.ID = gotResp.ID
				expResp.BundleID = gotResp.BundleID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[1].ID.String()),
			Token:      sd.Users[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &keyapp.NewKey{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"data\",\"error\":\"data is a required field\"},{\"field\":\"bundleID\",\"error\":\"bundleID is a required field\"},{\"field\":\"userID\",\"error\":\"userID is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"}]"),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[1].ID.String()),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[1].ID.String()),
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[1].ID.String()),
			Token:      sd.Users[0].Token + "A",
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[0].Bundles[2].ID.String()),
			Token:      sd.Users[1].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},

			ExpResp: errs.Newf(errs.Unauthenticated, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[1].ID.String(), sd.Users[0].Bundles[2].ID.String())
				return cmp.Diff(got, expResp)
			},
		},
	}

	return table
}
