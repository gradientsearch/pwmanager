package apitest

import (
	"net/http/httptest"
	"testing"

	authbuild "github.com/gradientsearch/pwmanager/api/services/auth/build/all"
	pwmanagerbuild "github.com/gradientsearch/pwmanager/api/services/pwmanager/build/all"
	"github.com/gradientsearch/pwmanager/app/sdk/auth"
	"github.com/gradientsearch/pwmanager/app/sdk/authclient"
	"github.com/gradientsearch/pwmanager/app/sdk/mux"
	"github.com/gradientsearch/pwmanager/business/sdk/dbtest"
)

// New initialized the system to run a test.
func New(t *testing.T, testName string) *Test {
	db := dbtest.New(t, testName)

	// -------------------------------------------------------------------------

	auth, err := auth.New(auth.Config{
		Log:       db.Log,
		DB:        db.DB,
		KeyLookup: &KeyStore{},
	})
	if err != nil {
		t.Fatal(err)
	}

	// -------------------------------------------------------------------------

	server := httptest.NewServer(mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus: db.BusDomain.User,
		},
		AuthConfig: mux.AuthConfig{
			Auth: auth,
		},
	}, authbuild.Routes()))

	authClient := authclient.New(db.Log, server.URL)

	// -------------------------------------------------------------------------

	mux := mux.WebAPI(mux.Config{
		Log: db.Log,
		DB:  db.DB,
		BusConfig: mux.BusConfig{
			UserBus:    db.BusDomain.User,
			KeyBus:     db.BusDomain.Key,
			BundleBus:  db.BusDomain.Bundle,
			VBundleBus: db.BusDomain.VBundle,
		},
		PwManagerConfig: mux.PwManagerConfig{
			AuthClient: authClient,
		},
	}, pwmanagerbuild.Routes())

	return &Test{
		DB:   db,
		Auth: auth,
		mux:  mux,
	}
}
