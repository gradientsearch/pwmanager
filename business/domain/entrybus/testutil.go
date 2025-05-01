package entrybus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/entry"
)

// TestGenerateNewEntries is a helper method for testing.
func TestGenerateNewEntries(n int, userID uuid.UUID, bids []uuid.UUID) []NewEntry {
	newEntries := make([]NewEntry, n)

	idx := rand.Intn(10000)
	for i := range n {
		idx++

		nk := NewEntry{
			Data:     entry.MustParse(fmt.Sprintf("Name%d", idx)),
			BundleID: bids[i],
			UserID:   userID,
		}

		newEntries[i] = nk
	}

	return newEntries
}

// TestGenerateSeedEntries is a helper method for testing.
func TestGenerateSeedEntries(ctx context.Context, n int, api *Business, userID uuid.UUID, bids []uuid.UUID) ([]Entry, error) {
	newEntries := TestGenerateNewEntries(n, userID, bids)

	entries := make([]Entry, len(newEntries))
	for i, nk := range newEntries {
		k, err := api.Create(ctx, nk)
		if err != nil {
			return nil, fmt.Errorf("seeding entry: idx: %d : %w", i, err)
		}

		entries[i] = k
	}

	return entries, nil
}
