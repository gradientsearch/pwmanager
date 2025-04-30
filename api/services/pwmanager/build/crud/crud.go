// Package crud binds the crud domain set of routes into the specified app.
package crud

import (
	"github.com/gradientsearch/pwmanager/app/domain/bundleapp"
	"github.com/gradientsearch/pwmanager/app/domain/checkapp"
	"github.com/gradientsearch/pwmanager/app/domain/keyapp"
	"github.com/gradientsearch/pwmanager/app/domain/tranapp"
	"github.com/gradientsearch/pwmanager/app/domain/userapp"
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
	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	bundleapp.Routes(app, bundleapp.Config{
		BundleBus:  cfg.BusConfig.BundleBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	keyapp.Routes(app, keyapp.Config{
		KeyBus:     cfg.BusConfig.KeyBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})

	tranapp.Routes(app, tranapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		KeyBus:     cfg.BusConfig.KeyBus,
		Log:        cfg.Log,
		AuthClient: cfg.PwManagerConfig.AuthClient,
		DB:         cfg.DB,
	})

	userapp.Routes(app, userapp.Config{
		UserBus:    cfg.BusConfig.UserBus,
		AuthClient: cfg.PwManagerConfig.AuthClient,
	})
}
