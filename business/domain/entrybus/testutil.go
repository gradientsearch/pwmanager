package entrybus

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/gradientsearch/pwmanager/business/types/entry"
)

const ()

// TestGenerateNewEntries is a helper method for testing.
// n specifies how many entries to add to each bundle.
// e.g. n = 10 than thn the the first 10 entries would be
// associated to the first bundle, the next 10 with the second bundle and so on.
func TestGenerateNewEntries(n int, userID uuid.UUID, bids []uuid.UUID) []NewEntry {
	newEntries := make([]NewEntry, n*len(bids))

	idx := rand.Intn(10000)

	for i := range len(bids) {
		idx++
		for b := range n {
			ne := NewEntry{
				Data:     entry.MustParse(fmt.Sprintf("Name%d", idx)),
				BundleID: bids[i],
				UserID:   userID,
			}

			newEntries[(i*n)+b] = ne
		}
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
