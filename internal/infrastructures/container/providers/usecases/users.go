// Package usecases wires application use cases.
package usecases

import (
	"github.com/darmayasa221/polymarket-go/internal/applications/security"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/commands/adduser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/getuser"
	"github.com/darmayasa221/polymarket-go/internal/applications/users/queries/listusers"
	userrepo "github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
)

// ProvideAddUser wires the AddUser use case.
func ProvideAddUser(repo userrepo.User, enc security.Encryption) adduser.UseCase {
	return adduser.New(repo, enc)
}

// ProvideGetUser wires the GetUser use case.
func ProvideGetUser(repo userrepo.User) getuser.UseCase {
	return getuser.New(repo)
}

// ProvideListUsers wires the ListUsers use case.
func ProvideListUsers(repo userrepo.User) listusers.UseCase {
	return listusers.New(repo)
}
