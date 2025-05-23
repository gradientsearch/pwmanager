package bundledb

import (
	"fmt"

	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
)

var orderByFields = map[string]string{
	bundlebus.OrderByID:     "bundle_id",
	bundlebus.OrderByType:   "type",
	bundlebus.OrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
