package entryapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	EntryBus   *entrybus.Business
	KeyBus     *keybus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAuthorizeEntry := mid.AuthorizeEntry(cfg.AuthClient, cfg.KeyBus, cfg.EntryBus)

	api := newApp(cfg.EntryBus)

	app.HandlerFunc(http.MethodGet, version, "bundles/{bundle_id}/entries/{entry_id}", api.queryByID, authen, ruleAuthorizeEntry)
	app.HandlerFunc(http.MethodPost, version, "bundles/{bundle_id}/entries", api.create, authen, ruleAuthorizeEntry)
	app.HandlerFunc(http.MethodPut, version, "bundles/{bundle_id}/entries/{entry_id}", api.update, authen, ruleAuthorizeEntry)
	app.HandlerFunc(http.MethodDelete, version, "bundles/{bundle_id}/entries/{entry_id}", api.delete, authen, ruleAuthorizeEntry)
}
