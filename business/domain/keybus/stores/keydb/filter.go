package keydb

import (
	"bytes"
	"strings"

	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

func (s *Store) applyFilter(filter keybus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["key_id"] = *filter.ID
		wc = append(wc, "key_id = :key_id")
	}

	if filter.UserID != nil {
		data["user_id"] = *filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
