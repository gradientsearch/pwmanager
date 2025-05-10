package uuk_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/sdk/uuk"
)

type TestUsersSecrets struct {
	Password  string
	SecretKey []byte
}

func TestUUK(t *testing.T) {
	uuk := uuk.UUK{}
	userSecrets := TestUsersSecrets{}
	userSecrets.Password = "gophers"
	userSecrets.SecretKey = make([]byte, 32)

	// For a random secret key use the following:
	//
	// if _, err := rand.Read(userSecrets.SecretKey); err != nil {
	// 	t.Fatalf("error creating secret key: %s; ", err)
	// }

	for i := range 32 {
		userSecrets.SecretKey[i] = byte('A')
	}

	groupID := uuid.New().String()
	userID := uuid.New().String()
	uuk.Build([]byte(userSecrets.Password), []byte(groupID), userSecrets.SecretKey, []byte(userID))
	fmt.Printf("%+v", uuk)
}
