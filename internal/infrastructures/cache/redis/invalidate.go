package redis

import "context"

// InvalidateUser removes a user from cache.
func (c *Client) InvalidateUser(ctx context.Context, userID string) error {
	return c.Delete(ctx, UserKey(userID))
}
