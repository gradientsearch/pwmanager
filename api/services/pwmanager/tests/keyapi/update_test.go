package key_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
)

func update200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       fmt.Sprintf("tu%d-%s-key", userBundleAdmin, userKeyMapping[userBundleAdmin]),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[userBundleAdmin].Keys[0].ID),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &keyapp.UpdateKey{
				Data: dbtest.StringPointer("Guitar"),
			},
			GotResp: &keyapp.Key{},
			ExpResp: &keyapp.Key{
				ID:          sd.Users[userBundleAdmin].Keys[0].ID.String(),
				UserID:      sd.Users[userBundleAdmin].ID.String(),
				BundleID:    sd.Users[userBundleAdmin].Bundles[0].ID.String(),
				Data:        "Guitar",
				Roles:       []string{"ADMIN", "READ", "WRITE"},
				DateCreated: sd.Users[userBundleAdmin].Keys[0].DateCreated.Format(time.RFC3339),
				DateUpdated: sd.Users[userBundleAdmin].Keys[0].DateCreated.Format(time.RFC3339),
			},
			CmpFunc: func(got any, exp any) string {
				gotResp, exists := got.(*keyapp.Key)
				if !exists {
					return "error occurred"
				}
				expResp := exp.(*keyapp.Key)
				gotResp.DateUpdated = expResp.DateUpdated

				return cmp.Diff(gotResp, expResp)
			},
		},
		{
			Name:       fmt.Sprintf("tu%d-%s-key-roles", userBundleAdmin, userKeyMapping[userBundleAdmin]),
			URL:        fmt.Sprintf("/v1/keys/role/%s", sd.Users[userBundleAdmin].Keys[0].ID),
			Token:      sd.Users[userBundleAdmin].Token,
			Method:     http.MethodPut,
			StatusCode: http.StatusOK,
			Input: &keyapp.Key{
				Roles: []string{"ADMIN", "READ", "WRITE"},
			},
			GotResp: &keyapp.Key{},
			ExpResp: &keyapp.UpdateBundleRole{},
			CmpFunc: func(got any, exp any) string {
				gotResp := got.(*keyapp.Key)
				gotRoles := keyapp.UpdateBundleRole{
					Roles: gotResp.Roles,
				}
				expRoles := keyapp.UpdateBundleRole{
					Roles: []string{
						bundlerole.Admin.String(),
						bundlerole.Read.String(),
						bundlerole.Write.String(),
					},
				}

				return cmp.Diff(gotRoles, expRoles)
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
			userBundleAdmin,
			"emptytoken",
			"&nbsp;",
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			userBundleAdmin,
			"badtoken",
			sd.Users[userBundleAdmin].Token[:10],
			errs.Newf(errs.Unauthenticated, "error parsing token: token contains an invalid number of segments"),
		},
		{
			userBundleAdmin,
			"badsig",
			sd.Users[userBundleAdmin].Token + "A",
			errs.Newf(errs.Unauthenticated, "authentication failed : bindings results[[{[true] map[x:false]}]] ok[true]"),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := apitest.Table{
			Name:       fmt.Sprintf("tu%d-%s", i.user, i.name),
			URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[i.user].Keys[0].ID),
			Token:      i.token,
			Method:     http.MethodPut,
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

func update403(sd apitest.SeedData) []apitest.Table {
	permError := fmt.Sprintf("must be an admin for bundle[%s] to modify a key", sd.Users[userBundleAdmin].Bundles[1].ID)

	inputs := []struct {
		user       userKey
		errMessage string
	}{
		{
			userReadWrite,
			permError,
		},
		{
			userRead,
			permError,
		},
		{
			userNoRoles,
			permError,
		},
		{
			userNoKey,
			fmt.Sprintf("query: userID[%s] bundleID[%s]: db: key not found", sd.Users[userNoKey].ID, sd.Users[userBundleAdmin].Bundles[1].ID),
		},
	}

	table := []apitest.Table{}
	for _, i := range inputs {
		t := []apitest.Table{
			{
				Name:       fmt.Sprintf("tu%d-%s-key", i.user, userKeyMapping[i.user]),
				URL:        fmt.Sprintf("/v1/keys/%s", sd.Users[userBundleAdmin].Keys[1].ID),
				Token:      sd.Users[i.user].Token,
				Method:     http.MethodPut,
				StatusCode: http.StatusForbidden,
				Input: &keyapp.UpdateKey{
					Data: dbtest.StringPointer("Guitar"),
				},
				GotResp: &errs.Error{},
				ExpResp: errs.Newf(errs.PermissionDenied, ""),
				CmpFunc: func(got any, exp any) string {
					expResp := exp.(*errs.Error)
					expResp.Message = i.errMessage
					return cmp.Diff(got, expResp)
				},
			},
			{
				Name:       fmt.Sprintf("tu%d-%s-role", i.user, userKeyMapping[i.user]),
				URL:        fmt.Sprintf("/v1/keys/role/%s", sd.Users[userBundleAdmin].Keys[1].ID),
				Token:      sd.Users[i.user].Token,
				Method:     http.MethodPut,
				StatusCode: http.StatusForbidden,
				Input: &keyapp.Key{
					Roles: []string{"ADMIN", "READ", "WRITE"},
				},
				GotResp: &errs.Error{},
				ExpResp: errs.Newf(errs.PermissionDenied, ""),
				CmpFunc: func(got any, exp any) string {
					expResp := exp.(*errs.Error)
					expResp.Message = i.errMessage
					return cmp.Diff(got, expResp)
				},
			},
		}
		table = append(table, t...)
	}

	return table
}
