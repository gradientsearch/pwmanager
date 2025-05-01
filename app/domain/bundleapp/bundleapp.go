// Package bundleapp maintains the app layer api for the bundle domain.
package bundleapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	bundleBus *bundlebus.Business
}

func newApp(bundleBus *bundlebus.Business) *app {
	return &app{
		bundleBus: bundleBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewBundle
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nb, err := toBusNewBundle(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	bdl, err := a.bundleBus.Create(ctx, nb)
	if err != nil {
		return errs.Newf(errs.Internal, "create: bdl[%+v]: %s", app, err)
	}

	return toAppBundle(bdl)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateBundle
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	uh, err := toBusUpdateBundle(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	bdl, err := mid.GetBundle(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "bundle missing in context: %s", err)
	}

	updUsr, err := a.bundleBus.Update(ctx, bdl, uh)
	if err != nil {
		return errs.Newf(errs.Internal, "update: bundleID[%s] uh[%+v]: %s", bdl.ID, uh, err)
	}

	return toAppBundle(updUsr)
}

func (a *app) delete(ctx context.Context, _ *http.Request) web.Encoder {
	bdl, err := mid.GetBundle(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "bundleID missing in context: %s", err)
	}

	if err := a.bundleBus.Delete(ctx, bdl); err != nil {
		return errs.Newf(errs.Internal, "delete: bundleID[%s]: %s", bdl.ID, err)
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

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, bundlebus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	bdls, err := a.bundleBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.bundleBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppBundles(bdls), total, page)
}

func (a *app) queryByID(ctx context.Context, _ *http.Request) web.Encoder {
	bdl, err := mid.GetBundle(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppBundle(bdl)
}
