package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/domain/bundlebus"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// AuthorizeEntryRetrieve validates a user has read permissions for a bundle entry.
func AuthorizeEntryRetrieve(client *authclient.Client, keyBus *keybus.Business, entryBus *entrybus.Business, bundleBus *bundlebus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// -------------------------------------------------------------------------
			// Validate

			entryID := web.Param(r, "entry_id")
			eID, err := uuid.Parse(entryID)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			// -------------------------------------------------------------------------
			// Get Entry

			entry, err := entryBus.QueryByID(ctx, eID)
			if err != nil {
				switch {
				case errors.Is(err, entrybus.ErrNotFound):
					return errs.New(errs.Unauthenticated, err)
				default:
					return errs.Newf(errs.Internal, "querybyid: entryID[%s] : %s", eID, err)
				}
			}

			// -------------------------------------------------------------------------
			// Authorize

			k, err := keyBus.QueryByUserIDBundleID(ctx, userID, entry.BundleID)
			if err != nil {
				switch {
				case errors.Is(err, keybus.ErrNotFound):
					return errs.New(errs.PermissionDenied, err)
				default:
					return errs.Newf(errs.Internal, "querybyid: userID[%s] bundleID[%s]: %s", userID, entry.BundleID, err)
				}
			}

			canRead := false
			for _, r := range k.Roles {
				if r.Equal(bundlerole.Read) {
					canRead = true
					break
				}
			}
			if !canRead {
				return errs.New(errs.PermissionDenied, fmt.Errorf("must have read perms for bundle[%s] to read entry", k.BundleID.String()))
			}

			// -------------------------------------------------------------------------
			// Set Entry

			ctx = setEntry(ctx, entry)

			return next(ctx, r)
		}

		return h
	}

	return m
}
