package entrydb

import (
	"fmt"

	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
)

var orderByFields = map[string]string{
	entrybus.OrderByEntryID: "entry_id",
	entrybus.OrderByUserID:  "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
