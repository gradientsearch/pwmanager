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
	keys := toAppVBundles(sd.Admins[0].User, sd.Admins[0].Keys)
	keys = append(keys, toAppVBundles(sd.Users[0].User, sd.Users[0].Keys)...)

	sort.Slice(keys, func(i, j int) bool {
		return keys[i].ID <= keys[j].ID
	})

	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/vbundles?page=1&rows=10&orderBy=key_id,ASC",
			Token:      sd.Admins[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &query.Result[vbundleapp.Key]{},
			ExpResp: &query.Result[vbundleapp.Key]{
				Page:        1,
				RowsPerPage: 10,
				Total:       len(keys),
				Items:       keys,
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
