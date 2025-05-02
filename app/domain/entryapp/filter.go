package entryapp

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/app/sdk/errs"
	"github.com/gradientsearch/pwmanager/business/domain/entrybus"
)

type queryParams struct {
	Page    string
	Rows    string
	OrderBy string
	ID      string
	UserID  string
}

func parseQueryParams(r *http.Request) queryParams {
	values := r.URL.Query()

	filter := queryParams{
		Page:    values.Get("page"),
		Rows:    values.Get("rows"),
		OrderBy: values.Get("orderBy"),
		ID:      values.Get("entry_id"),
		UserID:  values.Get("user_id"),
	}

	return filter
}

func parseFilter(qp queryParams) (entrybus.QueryFilter, error) {
	var fieldErrors errs.FieldErrors
	var filter entrybus.QueryFilter

	if qp.ID != "" {
		id, err := uuid.Parse(qp.ID)
		switch err {
		case nil:
			filter.ID = &id
		default:
			fieldErrors.Add("entry_id", err)
		}
	}

	if qp.UserID != "" {
		userID, err := uuid.Parse(qp.UserID)
		switch err {
		case nil:
			filter.UserID = &userID
		default:
			fieldErrors.Add("user_id", err)
		}
	}

	if fieldErrors != nil {
		return entrybus.QueryFilter{}, fieldErrors.ToError()
	}

	return filter, nil
}
