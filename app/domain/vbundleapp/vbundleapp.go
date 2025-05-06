// Package vbundleapp maintains the app layer api for the vbundle domain.
package vbundleapp

import (
	"context"
	"net/http"

	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/app/sdk/mid"
	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
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
	userID, err := mid.GetUserID(ctx)
	if err != nil {
		return errs.Newf(errs.InvalidArgument, "userID not found: %s", err)
	}

	bdls, err := a.vbundleBus.QueryByID(ctx, userID)
	if err != nil {
		return errs.Newf(errs.Internal, "query: %s", err)
	}

	return toAppUserBundleKeys(bdls)
}
