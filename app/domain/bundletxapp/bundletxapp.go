// Package bundletxapp maintains the app layer api for the tran domain.
package bundletxapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	userBus   *userbus.Business
	keyBus    *keybus.Business
	bundleBus *bundlebus.Business
}

func newApp(userBus *userbus.Business, keyBus *keybus.Business, bundleBus *bundlebus.Business) *app {
	return &app{
		userBus:   userBus,
		keyBus:    keyBus,
		bundleBus: bundleBus,
	}
}

// newWithTx constructs a new Handlers value with the domain apis
// using a store transaction that was created via middleware.
func (a *app) newWithTx(ctx context.Context) (*app, error) {
	tx, err := mid.GetTran(ctx)
	if err != nil {
		return nil, err
	}

	userBus, err := a.userBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	keyBus, err := a.keyBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	bundleBus, err := a.bundleBus.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	app := app{
		userBus:   userBus,
		keyBus:    keyBus,
		bundleBus: bundleBus,
	}

	return &app, nil
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewBundleTx
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nb, err := toBusNewBundle(ctx, app.Bundle)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	b, err := a.bundleBus.Create(ctx, nb)
	if err != nil {
		return errs.Newf(errs.Internal, "create: bdl[%+v]: %s", b, err)
	}

	nk, err := toBusNewKey(ctx, app.Key, b.ID)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	k, err := a.keyBus.Create(ctx, nk)
	if err != nil {
		return errs.Newf(errs.Internal, "create: key[%+v]: %s", k, err)
	}

	return toAppBundleTx(b, k)
}
