// Package keyapp maintains the app layer api for the key domain.
package keyapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	keyBus *keybus.Business
}

func newApp(keyBus *keybus.Business) *app {
	return &app{
		keyBus: keyBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewKey
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	np, err := toBusNewKey(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := a.keyBus.Create(ctx, np)
	if err != nil {
		return errs.Newf(errs.Internal, "create: prd[%+v]: %s", prd, err)
	}

	return toAppKey(prd)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateKey
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	up, err := toBusUpdateKey(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	prd, err := mid.GetKey(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "key missing in context: %s", err)
	}

	updPrd, err := a.keyBus.Update(ctx, prd, up)
	if err != nil {
		return errs.Newf(errs.Internal, "update: keyID[%s] up[%+v]: %s", prd.ID, app, err)
	}

	return toAppKey(updPrd)
}

func (a *app) delete(ctx context.Context, _ *http.Request) web.Encoder {
	prd, err := mid.GetKey(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "keyID missing in context: %s", err)
	}

	if err := a.keyBus.Delete(ctx, prd); err != nil {
		return errs.Newf(errs.Internal, "delete: keyID[%s]: %s", prd.ID, err)
	}

	return nil
}

func (a *app) query(ctx context.Context, r *http.Request) web.Encoder {
	qp := parseQueryParams(r)

	page, err := page.Parse(qp.Page, qp.Rows)
	if err != nil {
		return errs.NewFieldErrors("page", err)
	}

	filter, err := parseFilter(qp)
	if err != nil {
		return err.(*errs.Error)
	}

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, keybus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	keys, err := a.keyBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.keyBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppKeys(keys), total, page)
}

func (a *app) queryByID(ctx context.Context, r *http.Request) web.Encoder {
	prd, err := mid.GetKey(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppKey(prd)
}
