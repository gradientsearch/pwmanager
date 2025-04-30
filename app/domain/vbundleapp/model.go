package vbundleapp

import (
	"encoding/json"
	"time"

	"github.com/gradientsearch/pwmanager/business/domain/vbundlebus"
)

// Key represents information about an individual key with
// extended information.
type Key struct {
	ID          string `json:"id"`
	UserID      string `json:"userID"`
	Data        string `json:"data"`
	DateCreated string `json:"dateCreated"`
	DateUpdated string `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app Key) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppKey(prd vbundlebus.Key) Key {
	return Key{
		ID:          prd.ID.String(),
		UserID:      prd.UserID.String(),
		Data:        prd.Data.String(),
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppKeys(prds []vbundlebus.Key) []Key {
	app := make([]Key, len(prds))
	for i, prd := range prds {
		app[i] = toAppKey(prd)
	}

	return app
}
