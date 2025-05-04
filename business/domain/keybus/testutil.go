package keybus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/bundlerole"
	"github.com/gradientsearch/pwmanager/business/types/key"
)

// TestGenerateNewKeys is a helper method for testing.
func TestGenerateNewKeys(n int, userID uuid.UUID, bids []uuid.UUID, roles []bundlerole.Role) []NewKey {
	newKeys := make([]NewKey, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nk := NewKey{
			Data:     key.MustParse(fmt.Sprintf("Name%d", idx)),
			Roles:    roles,
			BundleID: bids[i],
			UserID:   userID,
		}

		newKeys[i] = nk
	}

	return newKeys
}

// TestGenerateSeedKeys is a helper method for testing.
func TestGenerateSeedKeys(ctx context.Context, n int, api *Business, userID uuid.UUID, bids []uuid.UUID, roles []bundlerole.Role) ([]Key, error) {
	newKeys := TestGenerateNewKeys(n, userID, bids, roles)

	keys := make([]Key, len(newKeys))
	for i, nk := range newKeys {
		k, err := api.Create(ctx, nk)
		if err != nil {
			return nil, fmt.Errorf("seeding key: idx: %d : %w", i, err)
		}

		keys[i] = k
	}

	return keys, nil
}
