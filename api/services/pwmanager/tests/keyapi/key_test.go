package key_test

import (
	"testing"

	"github.com/gradientsearch/pwmanager/app/sdk/apitest"
)

func Test_Key(t *testing.T) {
	t.Parallel()

	test := apitest.New(t, "Test_Key")

	// -------------------------------------------------------------------------

	sd, err := insertSeedData(test.DB, test.Auth)
	if err != nil {
		t.Fatalf("Seeding error: %s", err)
	}

	// -------------------------------------------------------------------------

	test.Run(t, queryByID200(sd), "querybyid-200")
	test.Run(t, queryByID401(sd), "querybyid-401")
	test.Run(t, queryByID403(sd), "querybyid-403")

	test.Run(t, create200(sd), "create-200")
	test.Run(t, create400(sd), "create-400")
	test.Run(t, create401(sd), "create-401")
	test.Run(t, create403(sd), "create-403")

	test.Run(t, update200(sd), "update-200")
	test.Run(t, update401(sd), "update-401")
	test.Run(t, update403(sd), "update-403")

	test.Run(t, delete401(sd), "delete-401")
	test.Run(t, delete403(sd), "delete-403")
	test.Run(t, delete200(sd), "delete-200")
}
