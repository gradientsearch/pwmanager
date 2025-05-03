package key_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[0].Keys[0].ID),
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &keyapp.Key{},
			ExpResp:    toAppKeyPtr(sd.Users[0].Keys[0]),
			CmpFunc: func(got any, exp any) string {

				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID401(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[0].Keys[0].ID),
			Token:      sd.Users[1].Token,
			StatusCode: http.StatusUnauthorized,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    &errs.Error{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*errs.Error)
				if !strings.Contains(gotResp.Message, "db: key not found") {
					return "message should have contained db: key not found"
				}
				return ""
			},
		},
	}

	return table
}
