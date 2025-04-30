package vbundle_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/business/domain/productbus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
)

func toAppVBundle(usr userbus.User, prd productbus.Product) vbundleapp.Product {
	return vbundleapp.Product{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Name:        prd.Name.String(),
		Cost:        prd.Cost.Value(),
		Quantity:    prd.Quantity.Value(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
		UserName:    usr.Name.String(),
	}
}

func toAppVBundles(usr userbus.User, prds []productbus.Product) []vbundleapp.Product {
	items := make([]vbundleapp.Product, len(prds))
	for i, prd := range prds {
		items[i] = toAppVBundle(usr, prd)
	}

	return items
}
