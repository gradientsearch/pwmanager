package entrybus

import "github.com/gradientsearch/pwmanager/business/sdk/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByEntryID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByEntryID = "entry_id"
	OrderByUserID  = "user_id"
)
