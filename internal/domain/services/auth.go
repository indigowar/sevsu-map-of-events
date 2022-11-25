package services

import (
	"context"
	"errors"
)

var (
	AuthErrFailedToLogin  = errors.New("failed to log in to account")
	AuthErrTokenIsInvalid = errors.New("token is invalid")
	AuthErrInternalError  = errors.New("internal error")
)

type AuthService interface {
	// Login - login as a user with given name and password
	// if the login process is failed return "", AuthErrFailedToLogin
	// if succeeded, then return RefreshToken and nil
	Login(ctx context.Context, name, password string) (string, error)

	// GetRefresh - when user wants to update their refresh token on a new one
	// if the token is invalid, it returns "", AuthErrTokenIsInvalid
	GetRefresh(ctx context.Context, rt string) (string, error)

	// GetAccess - returns an access token for user with refresh token = rt
	// if token is invalid, it returns "", AuthErrTokenIsInvalid
	// If some different error has happened, then it returns AuthErrInternalError
	GetAccess(ctx context.Context, rt string) (string, error)

	// CreateUser - creates a user with given name and password and returns a refresh token for this user
	CreateUser(ctx context.Context, name, password string) (string, error)
}
