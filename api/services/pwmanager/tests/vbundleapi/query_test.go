package vbundle_test

import (
	"net/http"

	"github.com/google/go-cmp/cmp"
	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
)

func query200(sd apitest.SeedData) []apitest.Table {
	table := []apitest.Table{
		{
			Name:       "basic",
			URL:        "/v1/vbundles",
			Token:      sd.Users[0].Token,
			StatusCode: http.StatusOK,
			Method:     http.MethodGet,
			GotResp:    &vbundleapp.UserBundleKeys{},
			ExpResp:    &vbundleapp.UserBundleKeys{},
			CmpFunc: func(got any, exp any) string {
				gotResp := *(got.(*vbundleapp.UserBundleKeys))
				if len(gotResp) != 2 {
					return "should have returned 2 bundles"
				}

				tu1 := sd.Users[0]
				tu2 := sd.Users[1]

				b1 := gotResp[0]
				b2 := gotResp[1]

				expResp := vbundleapp.UserBundleKeys{
					{
						UserID:      tu1.User.ID,
						BundleID:    tu1.Bundles[0].ID,
						Name:        tu1.User.Name.String(),
						Type:        b1.Type,
						Metadata:    tu1.Bundles[0].Metadata,
						DateCreated: b1.DateCreated,
						DateUpdated: b1.DateUpdated,
						KeyData:     tu1.Keys[0].Data.String(),
						KeyRoles:    bundlerole.ParseToString(tu1.Keys[0].Roles),
						Users: []vbundleapp.BundleUser{
							{
								UserID: tu1.User.ID,
								Name:   tu1.User.Name.String(),
								Email:  tu1.User.Email.Address,
								Roles:  bundlerole.ParseToString(tu1.Keys[0].Roles),
							},
							{
								UserID: tu2.User.ID,
								Name:   tu2.User.Name.String(),
								Email:  tu2.User.Email.Address,
								Roles:  bundlerole.ParseToString(tu2.Keys[0].Roles),
							},
						},
					},
					{
						UserID:      tu1.User.ID,
						BundleID:    tu1.Bundles[1].ID,
						Name:        tu1.User.Name.String(),
						Type:        b2.Type,
						Metadata:    tu1.Bundles[1].Metadata,
						DateCreated: b2.DateCreated,
						DateUpdated: b2.DateUpdated,
						KeyData:     tu1.Keys[1].Data.String(),
						KeyRoles:    bundlerole.ParseToString(tu1.Keys[1].Roles),
						Users: []vbundleapp.BundleUser{
							{
								UserID: tu1.User.ID,
								Name:   tu1.User.Name.String(),
								Email:  tu1.User.Email.Address,
								Roles:  bundlerole.ParseToString(tu1.Keys[1].Roles),
							},
						},
					},
				}

				return cmp.Diff(gotResp, expResp)
			},
		},
	}

	return table
}
