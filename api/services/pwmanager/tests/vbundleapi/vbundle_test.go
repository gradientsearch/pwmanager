package vbundle_test

import (
	"testing"

	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
)

func Test_VBundle(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_VBundle")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, query200(sd), "query-200")
}
