package key_test

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
			userBundleAdmin,
		},
	}

	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[i.user].Keys[0].ID),
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
			sd.Users[userBundleAdmin].Token[:10],
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			"badsig",
			sd.Users[userBundleAdmin].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	for _, i := range inputs {
		t := apitest.Table{

			Name:       fmt.Sprintf("tu%d-%s", userBundleAdmin, i.name),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[userBundleAdmin].Keys[2].ID),
			Token:      sd.Users[userBundleAdmin].Token + "A",
			Method:     http.MethodDelete,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
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
		errMessage string
	}{
		{
			userReadWrite,
			fmt.Sprintf("must be an admin for bundle[%s] to modify a key", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userRead,
			fmt.Sprintf("must be an admin for bundle[%s] to modify a key", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoRoles,
			fmt.Sprintf("must be an admin for bundle[%s] to modify a key", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoKey,
			fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[userBundleAdmin].Keys[0].ID),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodDelete,
			StatusCode: http.StatusForbidden,
			GotResp:    &errs.Error{},

			ExpResp: errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = i.errMessage
				return cmp.Diff(got, expResp)
			},
		}

		table = append(table, t)
	}

	return table
}
