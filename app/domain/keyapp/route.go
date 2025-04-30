package keyapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	KeyBus     *keybus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAny := mid.Authorize(cfg.AuthClient, auth.RuleAny)
	ruleUserOnly := mid.Authorize(cfg.AuthClient, auth.RuleUserOnly)
	ruleAuthorizeKey := mid.AuthorizeKey(cfg.AuthClient, cfg.KeyBus)

	api := newApp(cfg.KeyBus)

	app.HandlerFunc(http.MethodGet, version, "/keys", api.query, authen, ruleAny)
	app.HandlerFunc(http.MethodGet, version, "/keys/{key_id}", api.queryByID, authen, ruleAuthorizeKey)
	app.HandlerFunc(http.MethodPost, version, "/keys", api.create, authen, ruleUserOnly)
	app.HandlerFunc(http.MethodPut, version, "/keys/{key_id}", api.update, authen, ruleAuthorizeKey)
	app.HandlerFunc(http.MethodDelete, version, "/keys/{key_id}", api.delete, authen, ruleAuthorizeKey)
}
