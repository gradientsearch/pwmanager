package bundletx_test

import (
	"testing"

	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
)

func Test_BundleTx(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_BundleTx")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, create200(sd), "create-200")
	test.Run(t, create400(sd), "create-400")
}
