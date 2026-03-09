package redis

import "fmt"

// UserKey returns the cache key for a user by ID.
func UserKey(userID string) string { return fmt.Sprintf("user:%s", userID) }

// RateLimitKey returns the cache key for rate limiting.
func RateLimitKey(ip string) string { return fmt.Sprintf("rate_limit:%s", ip) }

// UserByUsernameKey returns the cache key for a user looked up by username.
func UserByUsernameKey(username string) string { return fmt.Sprintf("user:username:%s", username) }

// UserIDByUsernameKey returns the cache key for a user ID looked up by username.
func UserIDByUsernameKey(username string) string { return fmt.Sprintf("user:id:%s", username) }

// UserExistsKey returns the cache key for a username existence check.
func UserExistsKey(username string) string { return fmt.Sprintf("user:exists:%s", username) }

// AuthTokenKey returns the cache key for an authentication token by its value.
func AuthTokenKey(value string) string { return fmt.Sprintf("auth:token:%s", value) }

// AuthUserSetKey returns the cache key for the set of token keys belonging to a user.
func AuthUserSetKey(userID string) string { return fmt.Sprintf("auth:user:%s:tokens", userID) }
