// Package all binds all the routes into the specified app.
package all

import (
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/domain/bundletxapp"
	"github.com/gradientsearch/pwmanager/app/domain/checkapp"
	"github.com/gradientsearch/pwmanager/app/domain/entryapp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/domain/rawapp"
	"github.com/gradientsearch/pwmanager/app/domain/userapp"
	"github.com/gradientsearch/pwmanager/app/domain/vbundleapp"
	"github.com/gradientsearch/pwmanager/app/sdk/mux"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {
	// -------------------------------------------------------------------------
	// Service
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	rawapp.Routes(app)

	// -------------------------------------------------------------------------
	// Domains

	userapp.Routes(app, userapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	bundleapp.Routes(app, bundleapp.Config{
		Log:        cfg.Log,
		BundleBus:  cfg.BusConfig.BundleBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	keyapp.Routes(app, keyapp.Config{
		Log:        cfg.Log,
		KeyBus:     cfg.BusConfig.KeyBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	entryapp.Routes(app, entryapp.Config{
		Log:        cfg.Log,
		KeyBus:     cfg.BusConfig.KeyBus,
		EntryBus:   cfg.BusConfig.EntryBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	// -------------------------------------------------------------------------
	// TX

	bundletxapp.Routes(app, bundletxapp.Config{
		Log:        cfg.Log,
		DB:         cfg.DB,
		UserBus:    cfg.BusConfig.UserBus,
		KeyBus:     cfg.BusConfig.KeyBus,
		BundleBus:  cfg.BusConfig.BundleBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	vbundleapp.Routes(app, vbundleapp.Config{
		Log:        cfg.Log,
		UserBus:    cfg.BusConfig.UserBus,
		VBundleBus: cfg.BusConfig.VBundleBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})
}
