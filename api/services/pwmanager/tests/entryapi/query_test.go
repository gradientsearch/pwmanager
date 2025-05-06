package entry_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func queryByID200(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user userKey
	}{
		{
			userBundleAdmin,
		},
		{
			userReadWrite,
		},
		{
			userRead,
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[i.user].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &entryapp.Entry{},
			ExpResp:    toAppEntryPtr(sd.Users[userBundleAdmin].Entries[0]),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		}
		table = append(table, t)
	}

	return table
}

func queryByID401(sd apitest.SeedData) []apitest.Table {
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
			Name:       fmt.Sprintf("tu%d-%s-%s", userBundleAdmin, userKeyMapping[userBundleAdmin], i.name),
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      i.token,
			StatusCode: http.StatusUnauthorized,
			Method:     http.MethodGet,
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

func queryByID403(sd apitest.SeedData) []apitest.Table {
	permError := fmt.Sprintf("must have read perms for bundle[%s] to read entry", sd.Users[userBundleAdmin].Bundles[0].ID)

	inputs := []struct {
		user       userKey
		errMessage string
	}{
		{
			userNoRoles,
			permError,
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
			URL:        fmt.Sprintf("/v1/bundles/%s/entries/%s", sd.Users[userBundleAdmin].Bundles[0].ID, sd.Users[userBundleAdmin].Entries[0].ID),
			Token:      sd.Users[i.user].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = i.errMessage
				return cmp.Diff(got, exp)
			},
		}
		table = append(table, t)
	}
	return table
}
