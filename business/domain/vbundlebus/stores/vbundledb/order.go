package vbundledb

import (
	"fmt"

	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
)

var orderByFields = map[string]string{
	vbundlebus.OrderByKeyID:    "key_id",
	vbundlebus.OrderByUserID:   "user_id",
	vbundlebus.OrderByName:     "name",
	vbundlebus.OrderByCost:     "cost",
	vbundlebus.OrderByQuantity: "quantity",
	vbundlebus.OrderByUserName: "user_name",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
