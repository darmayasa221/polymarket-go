package providers

import (
	authrepo "github.com/darmayasa221/polymarket-go/internal/domains/authentications/repository"
	userrepo "github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	redisclient "github.com/darmayasa221/polymarket-go/internal/infrastructures/cache/redis"
	sqlitedb "github.com/darmayasa221/polymarket-go/internal/infrastructures/databases/sqlite"
	authredis "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/authentications/redis"
	usercached "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/users/cached"
	usersqlite "github.com/darmayasa221/polymarket-go/internal/infrastructures/repositories/users/sqlite"
)

// ProvideUserRepository wires the cached → sqlite user repository chain (Decorator Pattern).
func ProvideUserRepository(db *sqlitedb.DB, cache *redisclient.Client) userrepo.User {
	base := usersqlite.New(db)
	return usercached.New(base, cache)
}

// ProvideAuthRepository wires the Redis authentication repository.
func ProvideAuthRepository(cache *redisclient.Client) authrepo.Authentication {
	return authredis.New(cache)
}
