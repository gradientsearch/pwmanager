package keyapp

import (
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
)

var orderByFields = map[string]string{
	"key_id":   keybus.OrderByKeyID,
	"name":     keybus.OrderByName,
	"cost":     keybus.OrderByCost,
	"quantity": keybus.OrderByQuantity,
	"user_id":  keybus.OrderByUserID,
}
