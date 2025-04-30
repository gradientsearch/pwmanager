package bundleapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	BundleBus  *bundlebus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAny := mid.Authorize(cfg.AuthClient, auth.RuleAny)
	ruleUserOnly := mid.Authorize(cfg.AuthClient, auth.RuleUserOnly)
	ruleAuthorizeBundle := mid.AuthorizeBundle(cfg.AuthClient, cfg.BundleBus)

	api := newApp(cfg.BundleBus)

	app.HandlerFunc(http.MethodGet, version, "/bundles", api.query, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/bundles/{bundle_id}", api.queryByID, authen, ruleAuthorizeBundle)
	app.HandlerFunc(http.MethodPost, version, "/bundles", api.create, authen, ruleUserOnly)
	app.HandlerFunc(http.MethodPut, version, "/bundles/{bundle_id}", api.update, authen, ruleAuthorizeBundle)
	app.HandlerFunc(http.MethodDelete, version, "/bundles/{bundle_id}", api.delete, authen, ruleAuthorizeBundle)
}
