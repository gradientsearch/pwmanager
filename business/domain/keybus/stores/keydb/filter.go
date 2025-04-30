package keydb

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

func (s *Store) applyFilter(filter keybus.QueryFilter, data map[string]any, buf *bytes.Buffer) {
	var wc []string

	if filter.ID != nil {
		data["key_id"] = *filter.ID
		wc = append(wc, "key_id = :key_id")
	}

	if filter.Name != nil {
		data["name"] = fmt.Sprintf("%%%s%%", *filter.Name)
		wc = append(wc, "name LIKE :name")
	}

	if filter.Cost != nil {
		data["cost"] = *filter.Cost
		wc = append(wc, "cost = :cost")
	}

	if filter.Quantity != nil {
		data["quantity"] = *filter.Quantity
		wc = append(wc, "quantity = :quantity")
	}

	if len(wc) > 0 {
		buf.WriteString(" WHERE ")
		buf.WriteString(strings.Join(wc, " AND "))
	}
}
