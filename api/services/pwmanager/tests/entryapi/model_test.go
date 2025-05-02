package entry_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

func toAppKey(k keybus.Key) keyapp.Key {
	return keyapp.Key{
		ID:          k.ID.String(),
		UserID:      k.UserID.String(),
		Data:        k.Data.String(),
		DateCreated: k.DateCreated.Format(time.RFC3339),
		DateUpdated: k.DateUpdated.Format(time.RFC3339),
	}
}

func toAppKeyPtr(k keybus.Key) *keyapp.Key {
	appKey := toAppKey(k)
	return &appKey
}

func toAppKeys(keys []keybus.Key) []keyapp.Key {
	items := make([]keyapp.Key, len(keys))
	for i, k := range keys {
		items[i] = toAppKey(k)
	}

	return items
}
