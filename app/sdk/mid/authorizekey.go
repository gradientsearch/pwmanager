package mid

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// AuthorizeKeyRetrieve authorizes a user to retrieve their keys by key_id.
func AuthorizeKeyRetrieve(client *authclient.Client, keyBus *keybus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// -------------------------------------------------------------------------
			// Validation

			id := web.Param(r, "key_id")
			if id == "" {
				return errs.New(errs.InvalidArgument, ErrInvalidID)
			}

			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			keyID, err := uuid.Parse(id)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			// -------------------------------------------------------------------------
			// Get Key

			k, err := keyBus.QueryByID(ctx, keyID)
			if err != nil {
				switch {
				case errors.Is(err, keybus.ErrNotFound):
					return errs.New(errs.Unauthenticated, err)
				default:
					return errs.Newf(errs.Internal, "querybyid: keyID[%s]: %s", keyID, err)
				}
			}

			// -------------------------------------------------------------------------
			// Authorize

			if k.UserID != userID {
				return errs.New(errs.PermissionDenied, fmt.Errorf("only users can retrieve their own keys keyid[%s]", k.ID.String()))
			}

			// -------------------------------------------------------------------------
			// Set key

			ctx = setKey(ctx, k)

			return next(ctx, r)
		}

		return h
	}

	return m
}

// AuthorizeKeyCreate executes the specified role and extracts the specified
// key from the DB if a key id is specified in the call. Only Admins of a bundle can add new
// keys for users for that bundle.
func AuthorizeKeyCreate(client *authclient.Client, keyBus *keybus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// -------------------------------------------------------------------------
			// Validation

			id := web.Param(r, "bundle_id")
			if id == "" {
				return errs.New(errs.InvalidArgument, ErrInvalidID)
			}

			bundleID, err := uuid.Parse(id)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			// -------------------------------------------------------------------------
			// Authorize

			var k keybus.Key
			k, err = keyBus.QueryByUserIDBundleID(ctx, userID, bundleID)
			if err != nil {
				switch {
				case errors.Is(err, keybus.ErrNotFound):
					return errs.New(errs.PermissionDenied, err)
				default:
					return errs.Newf(errs.Internal, "querybyuseridbundleid: user_id[%s]  bundle_id[%s]: %s", userID, bundleID, err)
				}
			}

			isAdmin := false
			for _, r := range k.Roles {
				if r.Equal(bundlerole.Admin) {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				return errs.New(errs.PermissionDenied, fmt.Errorf("must have admin perms for bundle[%s] to create a key", k.BundleID.String()))
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

// AuthorizeKeyModify modifies a key if the calling user is a bundle admin.
func AuthorizeKeyModify(client *authclient.Client, keyBus *keybus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			// -------------------------------------------------------------------------
			// Validation
			id := web.Param(r, "key_id")
			if id == "" {
				return errs.New(errs.InvalidArgument, ErrInvalidID)
			}

			keyID, err := uuid.Parse(id)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			callerUserID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			// -------------------------------------------------------------------------
			// Get Key

			k, err := keyBus.QueryByID(ctx, keyID)
			if err != nil {
				switch {
				case errors.Is(err, keybus.ErrNotFound):
					return errs.New(errs.Unauthenticated, err)
				default:
					return errs.Newf(errs.Internal, "querybyid: keyID[%s]: %s", keyID, err)
				}
			}

			// -------------------------------------------------------------------------
			// Authorize

			var callerKey keybus.Key
			if k.UserID != callerUserID {
				callerKey, err = keyBus.QueryByUserIDBundleID(ctx, callerUserID, k.BundleID)
				if err != nil {
					switch {
					case errors.Is(err, keybus.ErrNotFound):
						return errs.New(errs.PermissionDenied, err)
					default:
						return errs.Newf(errs.Internal, "querybyid: keyID[%s]: %s", keyID, err)
					}
				}
			} else {
				// already have the key. No need to hit the db again.
				callerKey = k
			}

			isAdmin := false
			for _, r := range callerKey.Roles {
				if r.Equal(bundlerole.Admin) {
					isAdmin = true
					break
				}
			}
			if !isAdmin {
				return errs.New(errs.PermissionDenied, fmt.Errorf("must be an admin for bundle[%s] to modify a key", k.BundleID.String()))
			}

			// -------------------------------------------------------------------------
			// Set key

			ctx = setKey(ctx, k)

			return next(ctx, r)
		}

		return h
	}

	return m
}
