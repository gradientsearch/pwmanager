package bundleapp

import (
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
)

var orderByFields = map[string]string{
	"bundle_id": bundlebus.OrderByID,
	"type":      bundlebus.OrderByType,
	"user_id":   bundlebus.OrderByUserID,
}
