package bundle_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
)

func toAppBundle(bdl bundlebus.Bundle) bundleapp.Bundle {
	return bundleapp.Bundle{
		ID:          bdl.ID.String(),
		UserID:      bdl.UserID.String(),
		Type:        bdl.Type.String(),
		DateCreated: bdl.DateCreated.Format(time.RFC3339),
		DateUpdated: bdl.DateUpdated.Format(time.RFC3339),
	}
}

func toAppBundles(bundles []bundlebus.Bundle) []bundleapp.Bundle {
	items := make([]bundleapp.Bundle, len(bundles))
	for i, bdl := range bundles {
		items[i] = toAppBundle(bdl)
	}

	return items
}

func toAppBundlePtr(bdl bundlebus.Bundle) *bundleapp.Bundle {
	appBdl := toAppBundle(bdl)
	return &appBdl
}
