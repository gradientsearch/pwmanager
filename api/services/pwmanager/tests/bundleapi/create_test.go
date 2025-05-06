package bundle_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func create200(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user userKey
	}{
		{userA},
		{userB},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        "/v1/bundles",
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &bundleapp.NewBundleTx{
				Key: bundleapp.NewKey{
					Data: "Guitar",
				},
				Bundle: bundleapp.NewBundle{
					Type:     "PERSONAL",
					Metadata: "Bundle Metadata",
				},
			},
			GotResp: &bundleapp.BundleTx{},
			ExpResp: &bundleapp.BundleTx{
				Key: bundleapp.Key{
					Data: "Guitar",
				},
				Bundle: bundleapp.Bundle{
					Type:     "PERSONAL",
					Metadata: "Bundle Metadata",
				},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundleapp.BundleTx)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundleapp.BundleTx)
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
		}

		table = append(table, t)
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
			sd.Users[userA].Token[:10],
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			"badsig",
			sd.Users[userA].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s-%s", userA, userKeyMapping[userA], i.name),
			URL:        "/v1/bundles",
			Token:      i.token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			Input: &bundleapp.NewBundleTx{
				Bundle: bundleapp.NewBundle{
					Type:     "SPACE",
					Metadata: "BUNDLE METADATA",
				},
				Key: bundleapp.NewKey{
					Data: "ENCRYPTED SYMMETRIC KEY",
				},
			},
			GotResp: &errs.Error{},
			ExpResp: i.err,
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		}

		table = append(table, t)
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-%s-bad-type", userA, userKeyMapping[userA]),
			URL:        "/v1/bundles",
			Token:      sd.Users[userA].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.NewBundleTx{
				Bundle: bundleapp.NewBundle{
					Type:     "SPACE",
					Metadata: "BUNDLE METADATA",
				},
				Key: bundleapp.NewKey{
					Data: "ENCRYPTED SYMMETRIC KEY",
				},
			},
			GotResp: &errs.Error{},
			ExpResp: errs.Newf(errs.InvalidArgument, "parse: invalid bundle type \"SPACE\""),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-%s-missing-input", userA, userKeyMapping[userA]),
			URL:        "/v1/bundles",
			Token:      sd.Users[userA].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.NewBundleTx{
				Bundle: bundleapp.NewBundle{
					Type: "PERSONAL",
				},
			},
			GotResp: &errs.Error{},

			ExpResp: errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"bundle\",\"error\":\"bundle is a required field\"},{\"field\":\"key\",\"error\":\"key is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				gotResp := exp.(*errs.Error)

				expResp.FuncName = gotResp.FuncName
				expResp.FileName = gotResp.FileName

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
