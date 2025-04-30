package keydb

import (
	"fmt"

	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
)

var orderByFields = map[string]string{
	keybus.OrderByKeyID:    "key_id",
	keybus.OrderByUserID:   "user_id",
	keybus.OrderByName:     "name",
	keybus.OrderByCost:     "cost",
	keybus.OrderByQuantity: "quantity",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
