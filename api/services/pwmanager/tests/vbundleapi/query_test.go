package vbundle_test

import (
	"net/http"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
)

func query200(sd apitest.SeedData) []apitest.Table {
	prds := toAppVBundles(sd.Admins[0].User, sd.Admins[0].Products)
	prds = append(prds, toAppVBundles(sd.Users[0].User, sd.Users[0].Products)...)

	sort.Slice(prds, func(i, j int) bool {
		return prds[i].ID <= prds[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/vbundles?page=1&rows=10&orderBy=product_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[vbundleapp.Product]{},
			ExpResp: &query.Result[vbundleapp.Product]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(prds),
				Items:       prds,
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
			URL:        "/v1/vbundles?page=1&rows=10&name=$#!",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"name\",\"error\":\"invalid name \\\"$#!\\\"\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
		{
			Name:       "bad-orderby-value",
			URL:        "/v1/vbundles?page=1&rows=10&orderBy=roduct_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusBadRequest,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "[{\"field\":\"order\",\"error\":\"unknown order: roduct_id\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}
