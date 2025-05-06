package vbundle_test

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
)

func query200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/vbundles",
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &vbundleapp.UserBundleKeys{},
			ExpResp:    &vbundleapp.UserBundleKeys{},
			CmpFunc: func(got any, exp any) string {
				gotResp := *(got.(*vbundleapp.UserBundleKeys))
				if len(gotResp) != 2 {
					return "should have returned 2 bundles"
				}

				b1 := gotResp[0]

				if len(b1.Users) != 2 {
					return "should have returned 2 users in bundle 1"
				}

				b2 := gotResp[1]

				if len(b2.Users) != 1 {
					return "should have returned 1 users in bundle 2"
				}

				return ""
			},
		},
	}

	return table
}
