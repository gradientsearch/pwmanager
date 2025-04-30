package bundle_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
)

func toAppBundle(hme bundlebus.Bundle) bundleapp.Bundle {
	return bundleapp.Bundle{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type.String(),
		Address: bundleapp.Address{
			Address1: hme.Address.Address1,
			Address2: hme.Address.Address2,
			ZipCode:  hme.Address.ZipCode,
			City:     hme.Address.City,
			State:    hme.Address.State,
			Country:  hme.Address.Country,
		},
		DateCreated: hme.DateCreated.Format(time.RFC3339),
		DateUpdated: hme.DateUpdated.Format(time.RFC3339),
	}
}

func toAppBundles(bundles []bundlebus.Bundle) []bundleapp.Bundle {
	items := make([]bundleapp.Bundle, len(bundles))
	for i, hme := range bundles {
		items[i] = toAppBundle(hme)
	}

	return items
}

func toAppBundlePtr(hme bundlebus.Bundle) *bundleapp.Bundle {
	appHme := toAppBundle(hme)
	return &appHme
}
