package bundle_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func delete200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{}
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

	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[i.user].Bundles[0].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusNoContent,
		}

		table = append(table, t)
	}

	return table
}

func delete401(sd apitest.SeedData) []apitest.Table {
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
			Name:       fmt.Sprintf("tu%d-%s", userA, i.name),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[userA].Bundles[1].ID),
			Token:      i.token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    i.err,
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		}

		table = append(table, t)
	}

	return table
}

func delete403(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user       userKey
		deleteUser userKey
		err        *errs.Error
	}{
		{
			userA,
			userB,
			errs.Newf(errs.PermissionDenied, "only bundle owner can modify bundleID[%s]", sd.Users[userB].Bundles[1].ID),
		},
		{
			userB,
			userA,
			errs.Newf(errs.PermissionDenied, "only bundle owner can modify bundleID[%s]", sd.Users[userA].Bundles[1].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s", sd.Users[i.deleteUser].Bundles[1].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},
			ExpResp:    i.err,
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		}

		table = append(table, t)
	}
	return table
}
