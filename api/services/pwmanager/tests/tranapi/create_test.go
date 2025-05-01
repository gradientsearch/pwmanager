package tran_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/bundletxapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func create200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/tranexample",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &bundletxapp.NewBundleTx{
				Key: bundletxapp.NewKey{
					Data: "Guitar",
				},
				User: bundletxapp.NewUser{
					Name:            "Bill Kennedy",
					Email:           "bill@ardanlabs.com",
					Roles:           []string{"ADMIN"},
					Department:      "ITO",
					Password:        "123",
					PasswordConfirm: "123",
				},
			},
			GotResp: &bundletxapp.Key{},
			ExpResp: &bundletxapp.Key{
				Data: "Guitar",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundletxapp.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundletxapp.Key)

				expResp.ID = gotResp.ID
				expResp.UserID = gotResp.UserID
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
			URL:        "/v1/tranexample",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &bundletxapp.NewBundleTx{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"cost\",\"error\":\"cost is a required field\"},{\"field\":\"quantity\",\"error\":\"quantity is a required field\"},{\"field\":\"name\",\"error\":\"name is a required field\"},{\"field\":\"email\",\"error\":\"email is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"},{\"field\":\"password\",\"error\":\"password is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-name",
			URL:        "/v1/tranexample",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &bundletxapp.NewBundleTx{
				Key: bundletxapp.NewKey{
					Data: "Gu",
				},
				User: bundletxapp.NewUser{
					Name:            "Bill Kennedy",
					Email:           "bill@ardanlabs.com",
					Roles:           []string{"ADMIN"},
					Department:      "ITO",
					Password:        "123",
					PasswordConfirm: "123",
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "parse: invalid name \"Gu\""),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
