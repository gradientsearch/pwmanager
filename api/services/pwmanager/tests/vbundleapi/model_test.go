package vbundle_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
)

func toAppVBundle(usr userbus.User, prd keybus.Key) vbundleapp.Key {
	return vbundleapp.Key{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Data:        prd.Data.String(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppVBundles(usr userbus.User, prds []keybus.Key) []vbundleapp.Key {
	items := make([]vbundleapp.Key, len(prds))
	for i, prd := range prds {
		items[i] = toAppVBundle(usr, prd)
	}

	return items
}
