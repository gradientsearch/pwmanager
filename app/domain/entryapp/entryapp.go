// Package entryapp maintains the app layer api for the entry domain.
package entryapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	entryBus *entrybus.Business
}

func newApp(entryBus *entrybus.Business) *app {
	return &app{
		entryBus: entryBus,
	}
}

func (a *app) create(ctx context.Context, r *http.Request) web.Encoder {
	var app NewEntry
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ne, err := toBusNewEntry(ctx, app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	e, err := a.entryBus.Create(ctx, ne)
	if err != nil {
		return errs.Newf(errs.Internal, "create: k[%+v]: %s", e, err)
	}

	return toAppEntry(e)
}

func (a *app) update(ctx context.Context, r *http.Request) web.Encoder {
	var app UpdateEntry
	if err := web.Decode(r, &app); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	ue, err := toBusUpdateEntry(app)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	e, err := mid.GetEntry(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "entry missing in context: %s", err)
	}

	updEntry, err := a.entryBus.Update(ctx, e, ue)
	if err != nil {
		return errs.Newf(errs.Internal, "update: entryID[%s] uk[%+v]: %s", e.ID, app, err)
	}

	return toAppEntry(updEntry)
}

func (a *app) delete(ctx context.Context, _ *http.Request) web.Encoder {
	e, err := mid.GetEntry(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "entryID missing in context: %s", err)
	}

	if err := a.entryBus.Delete(ctx, e); err != nil {
		return errs.Newf(errs.Internal, "delete: entryID[%s]: %s", e.ID, err)
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

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, entrybus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	entries, err := a.entryBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.entryBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppEntrys(entries), total, page)
}

func (a *app) queryByID(ctx context.Context, r *http.Request) web.Encoder {
	e, err := mid.GetEntry(ctx)
	if err != nil {
		return errs.Newf(errs.Internal, "querybyid: %s", err)
	}

	return toAppEntry(e)
}
