package bundleapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	DB         *sqlx.DB
	UserBus    *userbus.Business
	KeyBus     *keybus.Business
	BundleBus  *bundlebus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAuthorizeBundleRetrieve := mid.AuthorizeBundleRetrieve(cfg.AuthClient, cfg.BundleBus)
	ruleAuthorizeBundleModify := mid.AuthorizeBundleModify(cfg.AuthClient, cfg.BundleBus)

	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	api := newApp(cfg.UserBus, cfg.KeyBus, cfg.BundleBus)

	app.HandlerFunc(http.MethodGet, version, "/bundles/{bundle_id}", api.queryByID, authen, ruleAuthorizeBundleRetrieve)
	app.HandlerFunc(http.MethodPost, version, "/bundles", api.create, authen, transaction)
	app.HandlerFunc(http.MethodPut, version, "/bundles/{bundle_id}", api.update, authen, ruleAuthorizeBundleModify)
	app.HandlerFunc(http.MethodDelete, version, "/bundles/{bundle_id}", api.delete, authen, ruleAuthorizeBundleModify)
}
