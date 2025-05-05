package bundle_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
)

func update200(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user userKey
	}{
		{
			userA,
		},
		{
			userB,
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[i.user].Bundles[0].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("PERSONAL"),
			},
			GotResp: &bundleapp.Bundle{},
			ExpResp: &bundleapp.Bundle{
				ID:     sd.Users[i.user].Bundles[0].ID.String(),
				UserID: sd.Users[i.user].ID.String(),
				Type:   "PERSONAL",
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*bundleapp.Bundle)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*bundleapp.Bundle)
				gotResp.DateUpdated = expResp.DateUpdated
				gotResp.DateCreated = expResp.DateCreated

				return cmp.Diff(gotResp, expResp)
			},
		}
		table = append(table, t)
	}
	return table
}

func update400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "bad-type",
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[0].ID),
			Token:      sd.Users[userA].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusBadRequest,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("BAD TYPE"),
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
	inputs := []struct {
		user  userKey
		name  string
		token string
		err   *errs.Error
	}{
		{
			userA,
			"emptytoken",
			"&nbsp;",
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			userA,
			"badtoken",
			sd.Users[userA].Token[:10],
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			userA,
			"badsig",
			sd.Users[userA].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, i.name),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[i.user].Bundles[0].ID),
			Token:      i.token,
			Method:     http.MethodPut,
			StatusCode: http.StatusUnauthorized,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("PERSONAL"),
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

func update403(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user         userKey
		userToUpdate userKey
		err          *errs.Error
	}{
		{
			userA,
			userB,
			errs.Newf(errs.PermissionDenied, "only bundle owner can modify bundleID[%s]", sd.Users[userB].Bundles[0].ID),
		},
		{
			userB,
			userA,
			errs.Newf(errs.PermissionDenied, "only bundle owner can modify bundleID[%s]", sd.Users[userA].Bundles[0].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[i.userToUpdate].Bundles[0].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusForbidden,
			Input: &bundleapp.UpdateBundle{
				Type: dbtest.StringPointer("PERSONAL"),
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
