package bundledb

import (
	"bytes"
	"strings"

	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
)

func (s *Store) applyFilter(filter bundlebus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["bundle_id"] = *filter.ID
		wc = append(wc, "bundle_id = :bundle_id")
	}

	if filter.UserID != nil {
		data["user_id"] = *filter.UserID
		wc = append(wc, "user_id = :user_id")
	}

	if filter.Type != nil {
		data["type"] = filter.Type.String()
		wc = append(wc, "type = :type")
	}

	if filter.StartCreatedDate != nil {
		data["start_date_created"] = filter.StartCreatedDate.UTC()
		wc = append(wc, "date_created >= :start_date_created")
	}

	if filter.EndCreatedDate != nil {
		data["end_date_created"] = filter.EndCreatedDate.UTC()
		wc = append(wc, "date_created <= :end_date_created")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
