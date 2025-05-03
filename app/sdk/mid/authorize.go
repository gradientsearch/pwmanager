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
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
	"github.com/gradientsearch/pwmanager/business/domain/keybus"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

// ErrInvalidID represents a condition where the id is not a uuid.
var ErrInvalidID = errors.New("ID is not in its proper form")

// Authorize validates authorization via the auth service.
func Authorize(client *authclient.Client, rule string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			auth := authclient.Authorize{
				Claims: GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			if err := client.Authorize(ctx, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			return next(ctx, r)
		}

		return h
	}

	return m
}

// AuthorizeUser executes the specified role and extracts the specified
// user from the DB if a user id is specified in the call. Depending on the rule
// specified, the userid from the claims may be compared with the specified
// user id.
func AuthorizeUser(client *authclient.Client, userBus *userbus.Business, rule string) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			id := web.Param(r, "user_id")

			var userID uuid.UUID

			if id != "" {
				var err error
				userID, err = uuid.Parse(id)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				usr, err := userBus.QueryByID(ctx, userID)
				if err != nil {
					switch {
					case errors.Is(err, userbus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Unauthenticated, "querybyid: userID[%s]: %s", userID, err)
					}
				}

				ctx = setUser(ctx, usr)
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			auth := authclient.Authorize{
				Claims: GetClaims(ctx),
				UserID: userID,
				Rule:   rule,
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

// AuthorizeKey executes the specified role and extracts the specified
// key from the DB if a key id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the key.
func AuthorizeKey(client *authclient.Client, keyBus *keybus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			id := web.Param(r, "key_id")

			// -------------------------------------------------------------------------
			// Validation

			if id == "" {
				return errs.New(errs.InvalidArgument, ErrInvalidID)
			}

			callerUserID, err := GetUserID(ctx)
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

			var callerKey keybus.Key
			if k.UserID != callerUserID {
				callerKey, err = keyBus.QueryByUserIDBundleID(ctx, callerUserID, k.BundleID)
				if err != nil {
					switch {
					case errors.Is(err, keybus.ErrNotFound):
						return errs.New(errs.Unauthenticated, err)
					default:
						return errs.Newf(errs.Internal, "querybyid: keyID[%s]: %s", keyID, err)
					}
				}
			} else {
				callerKey = k
			}

			if r.Method == "" || r.Method == http.MethodGet {
				if k.UserID != callerKey.UserID {
					return errs.New(errs.PermissionDenied, err)
				}
			} else {
				isAdmin := false
				for _, r := range callerKey.Roles {
					if r.Equal(bundlerole.Admin) {
						isAdmin = true
						break
					}
				}
				if !isAdmin {
					return errs.New(errs.PermissionDenied, err)
				}
			}

			// -------------------------------------------------------------------------

			ctx = setKey(ctx, k)

			return next(ctx, r)
		}

		return h
	}

	return m
}

// AuthorizeEntry executes the specified role and extracts the specified
// key from the DB if a user_id and bundle_id is specified in the call.
func AuthorizeEntry(client *authclient.Client, keyBus *keybus.Business, entryBus *entrybus.Business, bundleBus *bundlebus.Business) web.MidFunc {
	m := func(next web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, r *http.Request) web.Encoder {
			bundleID := web.Param(r, "bundle_id")
			entryID := web.Param(r, "entry_id")

			userID, err := GetUserID(ctx)
			if err != nil {
				return errs.New(errs.Unauthenticated, ErrInvalidID)
			}

			auth := authclient.Authorize{
				UserID: userID,
				Claims: GetClaims(ctx),
				Rule:   auth.RuleUserOnly,
			}

			if err := client.Authorize(ctx, auth); err != nil {
				return errs.New(errs.Unauthenticated, err)
			}

			entry := entrybus.Entry{
				UserID: userID,
			}
			if bundleID != "" {
				var err error
				bID, err := uuid.Parse(bundleID)
				if err != nil {
					return errs.New(errs.Unauthenticated, ErrInvalidID)
				}

				entry.BundleID = bID

				_, err = keyBus.QueryByUserIDBundleID(ctx, userID, bID)
				if err != nil {
					switch {
					case errors.Is(err, keybus.ErrNotFound):
						return errs.New(errs.PermissionDenied, err)
					default:
						return errs.Newf(errs.Internal, "querybyid: userID[%s] bundleID[%s]: %s", userID, bID, err)
					}
				}

				// TODO when roles are added to key, authorize user here
				// Check the request READ ROLE.
				if r.Method == "" || r.Method == http.MethodGet {
					//Check READ ACCESS
				} else {
					// Check WRITE ACCESS
				}

				// set the entry
				if r.Method != http.MethodPost {
					eID, err := uuid.Parse(entryID)
					if err != nil {
						return errs.New(errs.Unauthenticated, ErrInvalidID)
					}

					ce, err := entryBus.QueryByID(ctx, eID)
					if err != nil {
						switch {
						case errors.Is(err, entrybus.ErrNotFound):
							return errs.New(errs.Unauthenticated, err)
						default:
							return errs.Newf(errs.Internal, "querybyid: entryID[%s] : %s", eID, err)
						}
					}

					entry = ce
				}

				// Set the bundle
				if r.Method != "" && r.Method != http.MethodGet {
					bdl, err := bundleBus.QueryByID(ctx, bID)
					if err != nil {
						switch {
						case errors.Is(err, bundlebus.ErrNotFound):
							return errs.New(errs.Unauthenticated, err)
						default:
							return errs.Newf(errs.Internal, "querybyid: bundleID[%s] : %s", bID, err)
						}
					}

					ctx = setBundle(ctx, bdl)
				}

				ctx = setEntry(ctx, entry)
			} else {
				return errs.Newf(errs.InvalidArgument, "bundle_id query param is required")
			}

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			return next(ctx, r)
		}

		return h
	}

	return m
}

// AuthorizeBundle executes the specified role and extracts the specified
// bundle from the DB if a bundle id is specified in the call. Depending on
// the rule specified, the userid from the claims may be compared with the
// specified user id from the bundle.
func AuthorizeBundle(client *authclient.Client, bundleBus *bundlebus.Business) web.MidFunc {
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
