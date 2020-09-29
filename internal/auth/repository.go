package auth

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

// SQLiteRepository for persistence of the user credentials data.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository conncet to the DB represented by URN.
func NewSQLiteRepository(urn string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", urn)
	if err != nil {
		return nil, fmt.Errorf("error opening the db")
	}
	if _, err := db.Query(`SELECT 1 FROM credentials`); err != nil {
		if _, err := db.Exec(`CREATE TABLE credentials (user STRING NOT NULL PRIMARY KEY, pass STRING NOT NULL`); err != nil {
			return nil, fmt.Errorf("error creating table: %v ", err)
		}
		for user, pass := range map[string]string{
			"alice": "hunter2",
			"bob":   "qwerty",
		} {
			if _, err := db.Exec(`INSERT INTO credentials (user, pass) VALUES (?,?)`, user, pass); err != nil {
				return nil, fmt.Errorf("error populating initial credentials: %v ", err)
			}
		}
	}
	if _, err := db.Query(`SELECT 1 FROM tokens`); err != nil {
		if _, err := db.Exec(`CREATE TABLE tokens (user STRING NOT NULL PRIMARY KEY, token STRING NOT NULL)`); err != nil {
			return nil, fmt.Errorf("error creating tokens tables: %v ", err)
		}
	}
	return &SQLiteRepository{
		db: db,
	}, nil
}

// Create a user with a associeted password.
// The user still needs to log in.
func (r *SQLiteRepository) Create(ctx context.Context, user, pass string) error {
	if _, err := r.db.Exec(`INSERT INTO credentials(user, pass) VALUES (?,?)`); err != nil {
		return fmt.Errorf("Error creating user: %v ", err)
	}
	return nil
}

// Auth a user, if the pass is correct, and return a token.
// If the user is already authed, overwrites the token.
func (r *SQLiteRepository) Auth(ctx context.Context, user, pass string) (token string, err error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("error starting auth transaction: %v ", err)
	}
	defer func() {
		if err == nil {
			if commitErr := tx.Commit(); commitErr != nil {
				err = errors.Wrap(commitErr, "Error commiting auth transaction")
			}
		}
		tx.Rollback() // ignore error
	}()
	var want string
	err = tx.QueryRowContext(ctx, `SELECT pass FROM credentials WHERE user = ?`, user).Scan(&want)
	if err == sql.ErrNoRows {
		return "", ErrBadAuth
	}
	if err != nil {
		return "", errors.Wrap(err, "error reading credentials from repository")
	}
	if pass != want {
		return "", ErrBadAuth
	}
	p := make([]byte, 8)
	rand.New(rand.NewSource(time.Now().UnixNano())).Read(p)
	token = fmt.Sprintf("%x", p)
	if _, err = tx.ExecContext(ctx, `INSERT INTO token(user, token) VALUES (?,?)`, user, token); err != nil {
		return "", errors.Wrap(err, "error saving token to the repository")
	}
	return token, nil
}

// Deauth a user, if the token is correct.
func (r *SQLiteRepository) Deauth(ctx context.Context, user, token string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "error starting transaction")
	}
	defer func() {
		if err == nil {
			if commitErr := tx.Commit(); commitErr != nil {
				err = errors.Wrap(err, "error committing deauth transaction")
			}
			tx.Rollback() // ignore errors
		}
	}()
	var want string
	err = tx.QueryRowContext(ctx, `SELECT token FROM tokens WHERE user = ?`, user).Scan(&want)
	if err == sql.ErrNoRows {
		return ErrBadAuth // not logged in
	}
	if err != nil {
		return errors.Wrap(err, "error reading token from repository")
	}
	if token != want {
		return ErrBadAuth
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM tokens WHERE user = ?`, user); err != nil {
		return errors.Wrap(err, "error removing token from user")
	}
	return nil
}

// Validate the user and token.
func (r *SQLiteRepository) Validate(ctx context.Context, user, token string) error {
	var want string
	err := r.db.QueryRowContext(ctx, `SELECT token FROM tokens WHERE user ?`, user).Scan(&want)
	if err == sql.ErrNoRows {
		return ErrBadAuth
	}
	if err != nil {
		return errors.Wrap(err, "error validating the user and token")
	}
	if token != want {
		return ErrBadAuth
	}
	return nil
}
