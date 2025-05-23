package entryapp

import (
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/sdk/sqldb"
	"github.com/gradientsearch/pwmanager/foundation/logger"
	"github.com/gradientsearch/pwmanager/foundation/web"
	"github.com/jmoiron/sqlx"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log        *logger.Logger
	DB         *sqlx.DB
	EntryBus   *entrybus.Business
	BundleBus  *bundlebus.Business
	KeyBus     *keybus.Business
	AuthClient *authclient.Client
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	authen := mid.Authenticate(cfg.AuthClient)
	ruleAuthorizeEntryCreate := mid.AuthorizeEntryCreate(cfg.AuthClient, cfg.KeyBus, cfg.EntryBus, cfg.BundleBus)
	ruleAuthorizeEntryRetrieve := mid.AuthorizeEntryRetrieve(cfg.AuthClient, cfg.KeyBus, cfg.EntryBus, cfg.BundleBus)
	ruleAuthorizeEntryModify := mid.AuthorizeEntryModify(cfg.AuthClient, cfg.KeyBus, cfg.EntryBus, cfg.BundleBus)
	transaction := mid.BeginCommitRollback(cfg.Log, sqldb.NewBeginner(cfg.DB))

	api := newApp(cfg.EntryBus, cfg.BundleBus)

	app.HandlerFunc(http.MethodPost, version, "/bundles/{bundle_id}/entries", api.create, authen, ruleAuthorizeEntryCreate, transaction)
	app.HandlerFunc(http.MethodGet, version, "/entries/{entry_id}", api.queryByID, authen, ruleAuthorizeEntryRetrieve)
	app.HandlerFunc(http.MethodPut, version, "/entries/{entry_id}", api.update, authen, ruleAuthorizeEntryModify, transaction)
	app.HandlerFunc(http.MethodDelete, version, "/entries/{entry_id}", api.delete, authen, ruleAuthorizeEntryModify, transaction)
}
