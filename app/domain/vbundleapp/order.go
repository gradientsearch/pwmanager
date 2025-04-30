package vbundleapp

import (
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
)

var orderByFields = map[string]string{
	"key_id":    vbundlebus.OrderByKeyID,
	"user_id":   vbundlebus.OrderByUserID,
	"name":      vbundlebus.OrderByName,
	"cost":      vbundlebus.OrderByCost,
	"quantity":  vbundlebus.OrderByQuantity,
	"user_name": vbundlebus.OrderByUserName,
}
