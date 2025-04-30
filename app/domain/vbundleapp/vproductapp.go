// Package vbundleapp maintains the app layer api for the vbundle domain.
package vbundleapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/query"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
	"github.com/gradientsearch/pwmanager/business/sdk/order"
	"github.com/gradientsearch/pwmanager/business/sdk/page"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

type app struct {
	vbundleBus *vbundlebus.Business
}

func newApp(vbundleBus *vbundlebus.Business) *app {
	return &app{
		vbundleBus: vbundleBus,
	}
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

	orderBy, err := order.Parse(orderByFields, qp.OrderBy, vbundlebus.DefaultOrderBy)
	if err != nil {
		return errs.NewFieldErrors("order", err)
	}

	prds, err := a.vbundleBus.Query(ctx, filter, orderBy, page)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	total, err := a.vbundleBus.Count(ctx, filter)
	if err != nil {
		return errs.Newf(errs.Internal, "count: %s", err)
	}

	return query.NewResult(toAppKeys(prds), total, page)
}
