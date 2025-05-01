package bundletxapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/auth"
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
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))
	ruleAuthorizeUser := mid.Authorize(cfg.AuthClient, auth.RuleAdminOrSubject)

	api := newApp(cfg.UserBus, cfg.KeyBus, cfg.BundleBus)

	app.HandlerFunc(http.MethodPost, version, "/bundles", api.create, authen, ruleAuthorizeUser, transaction)
}
