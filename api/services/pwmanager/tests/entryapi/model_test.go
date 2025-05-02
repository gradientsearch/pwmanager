package entry_test

import (
	"time"

	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
)

func toAppEntry(e entrybus.Entry) entryapp.Entry {
	return entryapp.Entry{
		ID:          e.ID.String(),
		UserID:      e.UserID.String(),
		BundleID:    e.BundleID.String(),
		Data:        e.Data.String(),
		DateCreated: e.DateCreated.Format(time.RFC3339),
		DateUpdated: e.DateUpdated.Format(time.RFC3339),
	}
}

func toAppEntryPtr(e entrybus.Entry) *entryapp.Entry {
	appEntry := toAppEntry(e)
	return &appEntry
}

func toAppEntries(entries []entrybus.Entry) []entryapp.Entry {
	items := make([]entryapp.Entry, len(entries))
	for i, e := range entries {
		items[i] = toAppEntry(e)
	}

	return items
}
