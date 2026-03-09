// Package errors provides utilities for inspecting and unwrapping domain errors.
package errors

import (
	"errors"

	"github.com/darmayasa221/polymarket-go/internal/commons/errors/types"
)

// As attempts to unwrap err into the target type T.
func As[T error](err error) (T, bool) {
	var target T

	ok := errors.As(err, &target)

	return target, ok
}

// IsNotFound returns true if err is a *types.NotFoundError.
func IsNotFound(err error) bool {
	var target *types.NotFoundError

	return errors.As(err, &target)
}

// IsAuthentication returns true if err is a *types.AuthenticationError.
func IsAuthentication(err error) bool {
	var target *types.AuthenticationError

	return errors.As(err, &target)
}

// IsConflict returns true if err is a *types.ConflictError.
func IsConflict(err error) bool {
	var target *types.ConflictError

	return errors.As(err, &target)
}

// IsInvariant returns true if err is a *types.InvariantError.
func IsInvariant(err error) bool {
	var target *types.InvariantError

	return errors.As(err, &target)
}

// IsInternalServer returns true if err is a *types.InternalServerError.
func IsInternalServer(err error) bool {
	var target *types.InternalServerError

	return errors.As(err, &target)
}

// HTTPStatusOf returns the HTTP status code for a domain error, or 500 if unknown.
func HTTPStatusOf(err error) int {
	type statusCoder interface{ GetHTTPStatus() int }

	var sc statusCoder
	if errors.As(err, &sc) {
		return sc.GetHTTPStatus()
	}

	return 500
}

// CodeOf returns the error code for a domain error, or "" if unknown.
func CodeOf(err error) string {
	type codeCoder interface{ GetCode() string }

	var cc codeCoder
	if errors.As(err, &cc) {
		return cc.GetCode()
	}

	return ""
}
