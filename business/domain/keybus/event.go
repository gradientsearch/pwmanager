package keybus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gradientsearch/pwmanager/business/domain/userbus"
	"github.com/gradientsearch/pwmanager/business/sdk/delegate"
)

// registerDelegateFunctions will register action functions with the delegate
// system. If the business was constructed for query only, there won't be a
// delegate provided.
func (b *Business) registerDelegateFunctions() {
	if b.delegate != nil {
		b.delegate.Register(userbus.DomainName, userbus.ActionUpdated, b.actionUserUpdated)
	}
}

// actionUserUpdated is executed by the user domain indirectly when a user is updated.
func (b *Business) actionUserUpdated(ctx context.Context, data delegate.Data) error {
	var params userbus.ActionUpdatedParms
	err := json.Unmarshal(data.RawParams, &params)
	if err != nil {
		return fmt.Errorf("expected an encoded %T: %w", params, err)
	}

	b.log.Info(ctx, "action-userupdate", "user_id", params.UserID, "enabled", params.Enabled)

	// Now we can see if this user has been disabled. If they have been, we will
	// want to disable to mark all these keys as deleted. Right now we don't
	// have support for this, but you can see how we can process the event.

	return nil
}
