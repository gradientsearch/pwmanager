package key_test

import (
	"fmt"
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
)

func create200(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user userKey
	}{
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[userBundleAdmin].Bundles[2].ID.String()),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusOK,
			Input: &keyapp.NewKey{
				BundleID: sd.Users[userBundleAdmin].Bundles[2].ID.String(),
				UserID:   sd.Users[i.user].ID.String(),
				Data:     "Guitar",
				Roles:    []string{"ADMIN", "WRITE", "READ"},
			},
			GotResp: &keyapp.Key{},
			ExpResp: &keyapp.Key{
				Data:     "Guitar",
				UserID:   sd.Users[i.user].ID.String(),
				BundleID: sd.Users[userBundleAdmin].Bundles[2].ID.String(),
				Roles:    []string{"ADMIN", "WRITE", "READ"},
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*keyapp.Key)
				if !exists {
					return "error occurred"
				}

				expResp := exp.(*keyapp.Key)

				expResp.ID = gotResp.ID
				expResp.BundleID = gotResp.BundleID
				expResp.DateCreated = gotResp.DateCreated
				expResp.DateUpdated = gotResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		}
		table = append(table, t)
	}

	return table
}

func create400(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-missing-input", userBundleAdmin),
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[userBundleAdmin].Bundles[1].ID.String()),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusBadRequest,
			Input:      &keyapp.NewKey{},
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.InvalidArgument, "validate: [{\"field\":\"data\",\"error\":\"data is a required field\"},{\"field\":\"bundleID\",\"error\":\"bundleID is a required field\"},{\"field\":\"userID\",\"error\":\"userID is a required field\"},{\"field\":\"roles\",\"error\":\"roles is a required field\"}]"),
			CmpFunc: func(got any, exp any) string {
				return cmp.Diff(got, exp)
			},
		},
	}

	return table
}

func create401(sd apitest.SeedData) []apitest.Table {
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
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      i.token,
			Method:     http.MethodPost,
			StatusCode: http.StatusUnauthorized,
			GotResp:    &errs.Error{},
			ExpResp:    errs.Newf(errs.Unauthenticated, ""),
			CmpFunc: func(got any, exp any) string {
				expResp := exp.(*errs.Error)
				expResp.Message = i.err.Message
				return cmp.Diff(got, expResp)
			},
		}

		table = append(table, t)
	}
	return table
}

func create403(sd apitest.SeedData) []apitest.Table {
	inputs := []struct {
		user       userKey
		errMessage string
	}{
		{
			userRead,
			fmt.Sprintf("must have admin perms for bundle[%s] to create a key", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoRoles,
			fmt.Sprintf("must have admin perms for bundle[%s] to create a key", sd.Users[userBundleAdmin].Bundles[0].ID),
		},
		{
			userNoKey,
			fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[0].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", userBundleAdmin, userKeyMapping[i.user]),
			URL:        fmt.Sprintf("/v1/bundles/%s/keys", sd.Users[userBundleAdmin].Bundles[0].ID.String()),
			Token:      sd.Users[i.user].Token,
			Method:     http.MethodPost,
			StatusCode: http.StatusForbidden,
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
