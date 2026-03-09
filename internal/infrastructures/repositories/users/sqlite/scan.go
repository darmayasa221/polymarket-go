package sqlite

import (
	"database/sql"
	"errors"
	"time"

	errtypes "github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/repository"
	"github.com/darmayasa221/polymarket-go/internal/domains/users/user"
)

// scanUser scans a single *sql.Row into a User entity.
// Returns NotFoundError on sql.ErrNoRows, InternalServerError on other scan errors.
func scanUser(row *sql.Row) (*user.User, error) {
	var (
		id             string
		username       string
		email          string
		hashedPassword string
		fullName       string
		createdAt      time.Time
		updatedAt      time.Time
	)
	err := row.Scan(&id, &username, &email, &hashedPassword, &fullName, &createdAt, &updatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errtypes.NewNotFoundError(repository.ErrUserNotFound)
		}
		return nil, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	return user.Reconstitute(user.ReconstitutedParams{
		ID:             id,
		Username:       username,
		Email:          email,
		HashedPassword: hashedPassword,
		FullName:       fullName,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}), nil
}

// scanUsers scans all rows from a query into a slice of User entities.
// Returns InternalServerError on scan or iteration failure.
func scanUsers(rows *sql.Rows) ([]*user.User, error) {
	var users []*user.User
	for rows.Next() {
		var (
			id             string
			username       string
			email          string
			hashedPassword string
			fullName       string
			createdAt      time.Time
			updatedAt      time.Time
		)
		if err := rows.Scan(&id, &username, &email, &hashedPassword, &fullName, &createdAt, &updatedAt); err != nil {
			return nil, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
		}
		users = append(users, user.Reconstitute(user.ReconstitutedParams{
			ID:             id,
			Username:       username,
			Email:          email,
			HashedPassword: hashedPassword,
			FullName:       fullName,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		}))
	}
	if err := rows.Err(); err != nil {
		return nil, errtypes.NewInternalServerError(repository.ErrUserGetFailed)
	}
	return users, nil
}
