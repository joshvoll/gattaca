package auth

import (
	"context"
	"errors"
)

// Service describes the expected behavior of the authentication service.
// Users can create accounts, log in, and log out; other services can
// validate user tokens (sessions) that they've received.
type Service interface {
	Signup(ctx context.Context, user, pass string) error
	Login(ctx context.Context, user, pass string) (token string, err error)
	Logout(ctx context.Context, user, token string) error
	Validate(ctx context.Context, user, token string) error
}

// ErrBadAuth is a returned when authentication fails for any reason
var ErrBadAuth = errors.New("bad auth")

// DefaultService provide authentication via repository (DB).
// it's very thin layer around the repository
type DefaultService struct {
	repo Repository
}

// NewDefaultService return a usable service, wrapping a repository
func NewDefaultService(repo Repository) *DefaultService {
	return &DefaultService{
		repo: repo,
	}
}

// Signup create a user with a given password
// the user still needs to login
func (s *DefaultService) Signup(ctx context.Context, user, pass string) error {
	return s.repo.Create(ctx, user, pass)
}

// Login logs the user in, if the pass is correct.
// The returned token should be passed to logout or Validate
func (s *DefaultService) Login(ctx context.Context, user, pass string) (token string, err error) {
	return s.repo.Auth(ctx, user, pass)
}

// Logout log the user out, if the token is valid
func (s *DefaultService) Logout(ctx context.Context, user, token string) error {
	return s.repo.Deauth(ctx, user, token)
}

// Validate return nil error if the user is logged in and
// provdes the correct token
func (s *DefaultService) Validate(ctx context.Context, user, token string) (err error) {
	return s.repo.Validate(ctx, user, token)
}

// Repository models the data access layer required by the auth service.
// It's very similar to the service interface, because authentication
// doesn't involve much business logic.
type Repository interface {
	Create(ctx context.Context, user, pass string) error
	Auth(ctx context.Context, user, pass string) (token string, err error)
	Deauth(ctx context.Context, user, token string) error
	Validate(ctx context.Context, user, token string) error
}
