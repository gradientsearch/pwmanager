// Package tranapp maintains the app layer api for the tran domain.
package tranapp

import (
	"context"
	"errors"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	userBus *userbus.Business
	keyBus  *keybus.Business
}

func newApp(userBus *userbus.Business, keyBus *keybus.Business) *app {
	return &app{
		userBus: userBus,
		keyBus:  keyBus,
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

	app := app{
		userBus: userBus,
		keyBus:  keyBus,
	}

	return &app, nil
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewTran
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	a, err := a.newWithTx(ctx)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	nk, err := toBusNewKey(app.Key)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nu, err := toBusNewUser(app.User)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.userBus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.New(errs.Aborted, userbus.ErrUniqueEmail)
		}
		return errs.Newf(errs.Internal, "create: usr[%+v]: %s", usr, err)
	}

	nk.UserID = usr.ID

	k, err := a.keyBus.Create(ctx, nk)
	if err != nil {
		return errs.Newf(errs.Internal, "create: k[%+v]: %s", k, err)
	}

	return toAppKey(k)
}
