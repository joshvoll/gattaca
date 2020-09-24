package dna

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// Service describe the expected behavior of the DNA sequence test
// User can add their dna and check if subsecuencias exitst
type Service interface {
	Add(ctx context.Context, user, token, sequence string) error
	Check(ctx context.Context, user, token, subsequence string) error
}

// DefaultService provides our DNS sequence logic
type DefaultService struct {
	repo  Repository
	valid Validator
}

// Repository is a client-side interface, which models
// the concrete e.g. SQLiteRepository
type Repository interface {
	Insert(ctx context.Context, user, sequence string) error
	Select(ctx context.Context, user string) (sequence string, err error)
}

// Validator is a client-side interface, which models
// the parts of the auth service that we use.
type Validator interface {
	Validate(ctx context.Context, user, token string) error
}

var (
	// ErrSubsequenceNotFound is returned by Check on a failure.
	ErrSubsequenceNotFound = errors.New("subsequence doesn't appear in the DNA sequence")
	// ErrBadAuth is returned if a user validation check fails.
	ErrBadAuth = errors.New("bad auth")
	// ErrInvalidSequence is returned if an invalid sequence is added.
	ErrInvalidSequence = errors.New("invalid DNA sequence")
)

// NewDefaultService is the constructor function for the services
func NewDefaultService(repo Repository, valid Validator) *DefaultService {
	return &DefaultService{
		repo:  repo,
		valid: valid,
	}
}

// Add add a user and their sequence of dna to the database
func (s *DefaultService) Add(ctx context.Context, user, token, sequence string) error {
	if err := s.valid.Validate(ctx, user, token); err != nil {
		return ErrBadAuth
	}
	if !validSequence(sequence) {
		return ErrInvalidSequence
	}
	if err := s.repo.Insert(ctx, user, sequence); err != nil {
		return fmt.Errorf("error adding user: %v", err)
	}
	return nil
}

// Check returns true if the given subsequence is present in the user's DNA.
func (s *DefaultService) Check(ctx context.Context, user, token, subsequence string) error {
	if err := s.valid.Validate(ctx, user, token); err != nil {
		return ErrBadAuth
	}
	sequence, err := s.repo.Select(ctx, user)
	if err != nil {
		return fmt.Errorf("error reading sequence of dna: %v ", err)
	}
	if !strings.Contains(sequence, subsequence) {
		return ErrSubsequenceNotFound
	}
	return nil
}

func validSequence(sequence string) bool {
	for _, r := range sequence {
		switch r {
		case 'g', 'a', 't', 'c':
			continue
		default:
			return false
		}
	}
	return true
}
