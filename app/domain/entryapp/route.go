package entryapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/auth"
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
	ruleAny := mid.Authorize(cfg.AuthClient, auth.RuleAny)
	ruleUserOnly := mid.Authorize(cfg.AuthClient, auth.RuleUserOnly)
	ruleAuthorizeEntry := mid.AuthorizeEntry(cfg.AuthClient, cfg.KeyBus)

	api := newApp(cfg.EntryBus)

	app.HandlerFunc(http.MethodGet, version, "/entries", api.query, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/entries/{entry_id}", api.queryByID, authen, ruleAuthorizeEntry)
	app.HandlerFunc(http.MethodPost, version, "/entries", api.create, authen, ruleUserOnly)
	app.HandlerFunc(http.MethodPut, version, "/entries/{entry_id}", api.update, authen, ruleAuthorizeEntry)
	app.HandlerFunc(http.MethodDelete, version, "/entries/{entry_id}", api.delete, authen, ruleAuthorizeEntry)
}
