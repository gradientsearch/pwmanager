package mid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"

	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

func AuthorizeBundleRetrieve(client *authclient.Client, bundleBus *bundlebus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			id := web.Param(r, "bundle_id")

			var userID uuid.UUID

			if id != "" {
				var err error
				bundleID, err := uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				bdl, err := bundleBus.QueryByID(ctx, bundleID)
				if err != nil {
					switch {
					case errors.Is(err, bundlebus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: bundleID[%s]: %s", bundleID, err)
					}
				}

				userID = bdl.UserID
				ctx = setBundle(ctx, bdl)
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			auth := authclient.Authorize{
				Claims: GetClaims(ctx),
				UserID: userID,
				Rule:   auth.RuleAdminOrSubject,
			}

			if err := client.Authorize(ctx, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

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
