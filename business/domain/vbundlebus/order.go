package vbundlebus

import "github.com/gradientsearch/pwmanager/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByKeyID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByKeyID    = "key_id"
	OrderByUserID   = "user_id"
	OrderByName     = "name"
	OrderByCost     = "cost"
	OrderByQuantity = "quantity"
	OrderByUserName = "user_name"
)
