package key_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
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
		{
			userNoRoles,
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[i.user].Keys[0].ID),
			Token:      sd.Users[i.user].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &keyapp.Key{},
			ExpResp:    toAppKeyPtr(sd.Users[i.user].Keys[0]),
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
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[0].Keys[0].ID),
			Token:      i.token,
			StatusCode: http.StatusUnauthorized,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    &errs.Error{},
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp = i.err
				return cmp.Diff(got, expResp)
			},
		}

		table = append(table, t)
	}

	return table
}

func queryByID403(sd apitest.SeedData) []apitest.Table {
	invalidKey := uuid.New()
	inputs := []struct {
		user       userKey
		keyID      uuid.UUID
		errMessage string
	}{
		{
			userBundleAdmin,
			sd.Users[userReadWrite].Keys[0].ID,
			fmt.Sprintf("only users can retrieve their own keys keyid[%s]", sd.Users[userReadWrite].Keys[0].ID),
		},
		{
			userNoRoles,
			sd.Users[userBundleAdmin].Keys[0].ID,
			fmt.Sprintf("only users can retrieve their own keys keyid[%s]", sd.Users[userBundleAdmin].Keys[0].ID),
		},
		{
			userNoKey,
			invalidKey,
			fmt.Sprintf("query: keyID[%s]: db: key not found", invalidKey),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/keys/%s", &i.keyID),
			Token:      sd.Users[i.user].Token,
			StatusCode: http.StatusForbidden,
			Method:     http.MethodGet,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.PermissionDenied, ""),
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
