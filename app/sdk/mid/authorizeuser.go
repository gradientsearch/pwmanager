package mid

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/foundation/web"
)

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
