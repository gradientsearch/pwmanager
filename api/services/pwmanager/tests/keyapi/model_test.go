package key_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

func toAppKey(prd keybus.Key) keyapp.Key {
	return keyapp.Key{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Data:        prd.Data.String(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppKeyPtr(prd keybus.Key) *keyapp.Key {
	appPrd := toAppKey(prd)
	return &appPrd
}

func toAppKeys(prds []keybus.Key) []keyapp.Key {
	items := make([]keyapp.Key, len(prds))
	for i, prd := range prds {
		items[i] = toAppKey(prd)
	}

	return items
}
