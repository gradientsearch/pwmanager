package mid

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"

	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// AuthorizeBundleModify validates the user is the bundle owner prior to modifying the bundle.
func AuthorizeBundleModify(client *authclient.Client, bundleBus *bundlebus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// -------------------------------------------------------------------------
			// Validate Input

			bid := web.Param(r, "bundle_id")
			bundleID, err := uuid.Parse(bid)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			// -------------------------------------------------------------------------
			// Get bundle

			bdl, err := bundleBus.QueryByID(ctx, bundleID)
			if err != nil {
				switch {
				case errors.Is(err, bundlebus.ErrNotFound):
					return errs.New(errs.PermissionDenied, err)
				default:
					return errs.Newf(errs.Unauthenticated, "querybyid: bundleID[%s]: %s", bundleID, err)
				}
			}

			// -------------------------------------------------------------------------
			// Authorize

			if userID != bdl.UserID {
				return errs.Newf(errs.PermissionDenied, "only bundle owner can modify bundleID[%s]", bdl.ID)
			}

			// -------------------------------------------------------------------------
			// Set bundle

			ctx = setBundle(ctx, bdl)

			return next(ctx, r)
		}
		return h
	}

	return m
}
