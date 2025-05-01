package vbundle_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
)

func toAppVBundle(usr userbus.User, k keybus.Key) vbundleapp.Key {
	return vbundleapp.Key{
		ID:          k.ID.String(),
		UserID:      k.UserID.String(),
		Data:        k.Data.String(),
		DateCreated: k.DateCreated.Format(time.RFC3339),
		DateUpdated: k.DateUpdated.Format(time.RFC3339),
	}
}

func toAppVBundles(usr userbus.User, keys []keybus.Key) []vbundleapp.Key {
	items := make([]vbundleapp.Key, len(keys))
	for i, k := range keys {
		items[i] = toAppVBundle(usr, k)
	}

	return items
}
