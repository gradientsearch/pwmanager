package bundlebus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundletype"
)

// TestGenerateNewBundles is a helper method for testing.
func TestGenerateNewBundles(n int, userID uuid.UUID) []NewBundle {
	newHmes := make([]NewBundle, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		t := bundletype.Single
		if v := (idx + i) % 2; v == 0 {
			t = bundletype.Condo
		}

		nh := NewBundle{
			Type: t,
			Address: Address{
				Address1: fmt.Sprintf("Address%d", idx),
				Address2: fmt.Sprintf("Address%d", idx),
				ZipCode:  fmt.Sprintf("%05d", idx),
				City:     fmt.Sprintf("City%d", idx),
				State:    fmt.Sprintf("State%d", idx),
				Country:  fmt.Sprintf("Country%d", idx),
			},
			UserID: userID,
		}

		newHmes[i] = nh
	}

	return newHmes
}

// TestGenerateSeedBundles is a helper method for testing.
func TestGenerateSeedBundles(ctx context.Context, n int, api *Business, userID uuid.UUID) ([]Bundle, error) {
	newHmes := TestGenerateNewBundles(n, userID)

	hmes := make([]Bundle, len(newHmes))
	for i, nh := range newHmes {
		hme, err := api.Create(ctx, nh)
		if err != nil {
			return nil, fmt.Errorf("seeding bundle: idx: %d : %w", i, err)
		}

		hmes[i] = hme
	}

	return hmes, nil
}

// ParseAddress is a helper function to create an address value.
func ParseAddress(address1 string, address2 string, zipCode string, city string, state string, country string) Address {
	return Address{
		Address1: address1,
		Address2: address2,
		ZipCode:  zipCode,
		City:     city,
		State:    state,
		Country:  country,
	}
}
