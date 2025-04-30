package dbtest

import (
	"time"

	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus/stores/bundledb"
	"github.com/gradientsearch/pwmanager/business/domain/productbus"
	"github.com/gradientsearch/pwmanager/business/domain/productbus/stores/productdb"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus/stores/usercache"
	"github.com/gradientsearch/pwmanager/business/domain/userbus/stores/userdb"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus/stores/vbundledb"
	"github.com/gradientsearch/pwmanager/business/sdk/delegate"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/jmoiron/sqlx"
)

// BusDomain represents all the business domain apis needed for testing.
type BusDomain struct {
	Delegate *delegate.Delegate
	Bundle   *bundlebus.Business
	Product  *productbus.Business
	User     *userbus.Business
	VBundle  *vbundlebus.Business
}

func newBusDomains(log *logger.Logger, db *sqlx.DB) BusDomain {
	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, usercache.NewStore(log, userdb.NewStore(log, db), time.Hour))
	productBus := productbus.NewBusiness(log, userBus, delegate, productdb.NewStore(log, db))
	bundleBus := bundlebus.NewBusiness(log, userBus, delegate, bundledb.NewStore(log, db))
	vbundleBus := vbundlebus.NewBusiness(vbundledb.NewStore(log, db))

	return BusDomain{
		Delegate: delegate,
		Bundle:   bundleBus,
		Product:  productBus,
		User:     userBus,
		VBundle:  vbundleBus,
	}
}
