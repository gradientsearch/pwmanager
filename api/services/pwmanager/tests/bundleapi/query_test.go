package bundle_test

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
)

func query200(sd apitest.SeedData) []apitest.Table {
	hmes := make([]bundlebus.Bundle, 0, len(sd.Admins[0].Bundles)+len(sd.Users[0].Bundles))
	hmes = append(hmes, sd.Admins[0].Bundles...)
	hmes = append(hmes, sd.Users[0].Bundles...)

	sort.Slice(hmes, func(i, j int) bool {
		return hmes[i].ID.String() <= hmes[j].ID.String()
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/bundles?page=1&rows=10&orderBy=bundle_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[bundleapp.Bundle]{},
			ExpResp: &query.Result[bundleapp.Bundle]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(hmes),
				Items:       toAppBundles(hmes),
			},
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func query400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-query-filter",
			URL:        "/v1/bundles?page=1&rows=10&type=bungalow",
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"type\",\"error\":\"invalid bundle type \\\"bungalow\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-orderby-value",
			URL:        "/v1/bundles?page=1&rows=10&orderBy=ome_id,ASC",
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"order\",\"error\":\"unknown order: ome_id\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func queryByID200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[0].Bundles[0].ID),
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &bundleapp.Bundle{},
			ExpResp:    toAppBundlePtr(sd.Users[0].Bundles[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
