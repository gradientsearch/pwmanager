package keyapp

import (
	"net/http"

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
	ruleAuthorizeKeyCreate := mid.AuthorizeKeyCreate(cfg.AuthClient, cfg.KeyBus)
	ruleAuthorizeKeyRetrieve := mid.AuthorizeKeyRetrieve(cfg.AuthClient, cfg.KeyBus)
	ruleAuthorizeKeyModify := mid.AuthorizeKeyModify(cfg.AuthClient, cfg.KeyBus)

	api := newApp(cfg.KeyBus)

	app.HandlerFunc(http.MethodPost, version, "/bundles/{bundle_id}/keys", api.create, authen, ruleAuthorizeKeyCreate)
	app.HandlerFunc(http.MethodGet, version, "/keys/{key_id}", api.queryByID, authen, ruleAuthorizeKeyRetrieve)
	app.HandlerFunc(http.MethodPut, version, "/keys/role/{key_id}", api.updateRole, authen, ruleAuthorizeKeyModify)
	app.HandlerFunc(http.MethodPut, version, "/keys/{key_id}", api.update, authen, ruleAuthorizeKeyModify)
	app.HandlerFunc(http.MethodDelete, version, "/keys/{key_id}", api.delete, authen, ruleAuthorizeKeyModify)
}
