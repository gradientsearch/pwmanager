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

		t := bundletype.Personal
		if v := (idx + i) % 2; v == 0 {
			t = bundletype.Shareable
		}

		nh := NewBundle{
			Type:   t,
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
