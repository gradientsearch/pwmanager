package entryapp

import (
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
)

var orderByFields = map[string]string{
	"entry_id": entrybus.OrderByEntryID,
	"user_id":  entrybus.OrderByUserID,
}
