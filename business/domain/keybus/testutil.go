package keybus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/key"
)

// TestGenerateNewKeys is a helper method for testing.
func TestGenerateNewKeys(n int, userID uuid.UUID) []NewKey {
	newPrds := make([]NewKey, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		np := NewKey{
			Data:   key.MustParse(fmt.Sprintf("Name%d", idx)),
			UserID: userID,
		}

		newPrds[i] = np
	}

	return newPrds
}

// TestGenerateSeedKeys is a helper method for testing.
func TestGenerateSeedKeys(ctx context.Context, n int, api *Business, userID uuid.UUID) ([]Key, error) {
	newPrds := TestGenerateNewKeys(n, userID)

	prds := make([]Key, len(newPrds))
	for i, np := range newPrds {
		prd, err := api.Create(ctx, np)
		if err != nil {
			return nil, fmt.Errorf("seeding key: idx: %d : %w", i, err)
		}

		prds[i] = prd
	}

	return prds, nil
}
