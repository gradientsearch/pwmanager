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
	newBdls := make([]NewBundle, n)

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

		newBdls[i] = nh
	}

	return newBdls
}

// TestGenerateSeedBundles is a helper method for testing.
func TestGenerateSeedBundles(ctx context.Context, n int, api *Business, userID uuid.UUID) ([]Bundle, error) {
	newBdls := TestGenerateNewBundles(n, userID)

	bdls := make([]Bundle, len(newBdls))
	for i, nh := range newBdls {
		bdl, err := api.Create(ctx, nh)
		if err != nil {
			return nil, fmt.Errorf("seeding bundle: idx: %d : %w", i, err)
		}

		bdls[i] = bdl
	}

	return bdls, nil
}
