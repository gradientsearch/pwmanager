package keyapp

import (
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

var orderByFields = map[string]string{
	"key_id":  keybus.OrderByKeyID,
	"user_id": keybus.OrderByUserID,
}
