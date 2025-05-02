package entry_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
)

func toAppEntry(k entrybus.Entry) entryapp.Entry {
	return entryapp.Entry{
		ID:          k.ID.String(),
		UserID:      k.UserID.String(),
		Data:        k.Data.String(),
		DateCreated: k.DateCreated.Format(time.RFC3339),
		DateUpdated: k.DateUpdated.Format(time.RFC3339),
	}
}

func toAppEntryPtr(k entrybus.Entry) *entryapp.Entry {
	appEntry := toAppEntry(k)
	return &appEntry
}

func toAppEntries(entries []entrybus.Entry) []entryapp.Entry {
	items := make([]entryapp.Entry, len(entries))
	for i, k := range entries {
		items[i] = toAppEntry(k)
	}

	return items
}
