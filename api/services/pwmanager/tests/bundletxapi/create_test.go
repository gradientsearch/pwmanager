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
			URL:        "/v1/bundles",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &bundletxapp.NewBundleTx{
				Key: bundletxapp.NewKey{
					Data: "Guitar",
				},
				Bundle: bundletxapp.NewBundle{
					Type:     "PERSONAL",
					Metadata: "Bundle Metadata",
				},
			},
			GotResp: &bundletxapp.BundleTx{},
			ExpResp: &bundletxapp.BundleTx{
				Key: bundletxapp.Key{
					Data: "Guitar",
				},
				Bundle: bundletxapp.Bundle{
					Type:     "PERSONAL",
					Metadata: "Bundle Metadata",
				},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundletxapp.BundleTx)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundletxapp.BundleTx)

				expResp.Bundle.ID = gotResp.Bundle.ID
				expResp.Bundle.DateCreated = gotResp.Bundle.DateCreated
				expResp.Bundle.DateUpdated = gotResp.Bundle.DateUpdated
				expResp.Bundle.UserID = gotResp.Bundle.UserID

				expResp.Key.ID = gotResp.Key.ID
				expResp.Key.DateCreated = gotResp.Key.DateCreated
				expResp.Key.DateUpdated = gotResp.Key.DateUpdated
				expResp.Key.UserID = gotResp.Key.UserID

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
			URL:        "/v1/bundles",
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
			URL:        "/v1/bundles",
			Token:      sd.Admins[0].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &bundletxapp.NewBundleTx{
				Key: bundletxapp.NewKey{
					Data: "Gu",
				},
				Bundle: bundletxapp.NewBundle{
					Type:     "PERSONAL",
					Metadata: "Bundle Metadata",
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
