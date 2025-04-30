package bundlebus

import "github.com/gradientsearch/pwmanager/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID     = "bundle_id"
	OrderByType   = "type"
	OrderByUserID = "user_id"
)
