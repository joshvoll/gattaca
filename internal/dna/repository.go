package dna

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // driver for sql lite
)

// ErrInvalidUser is returned when an invalid user is passed to Select.
var (
	ErrInvalidUser = errors.New("invalid user")
)

// SQLiteRepository for persistence of the DNA project.
type SQLiteRepository struct {
	db *sql.DB
}

// NewSQLiteRepository connects to the DB represented by URN.
// For now we're going to create the db on the same constructor function
func NewSQLiteRepository(urn string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", urn)
	if err != nil {
		return nil, fmt.Errorf("Error opening the sql lite db: %v ", err)
	}
	if _, err := db.Query(`SELECT 1 FROM dna`); err != nil {
		if _, err := db.Exec(`CREATE TABLE dna (user STRING NOT NULL PRIMARY KEY, sequence STRING NOT NULL)`); err != nil {
			return nil, fmt.Errorf("error creating table dna: %v ", err)
		}
		for user, sequence := range map[string]string{
			"alice": "attcgtattattttttgatatttttccacaaaaatacagactaaatacaactgaatacag",
			"bob":   "tgcaaaattagatataaatgtaaacgaacataaaaacttttataagacaggattaagtta",
		} {
			if _, err := db.Exec(`INSERT INTO dna (user, sequence) VALUES (?,?)`, user, sequence); err != nil {
				return nil, fmt.Errorf("error populating initial sequences: %v ", err)
			}
		}
	}
	return &SQLiteRepository{
		db: db,
	}, nil
}

// Insert inser a user with the dna sequence and return and error is ocurrs
func (s *SQLiteRepository) Insert(ctx context.Context, user, sequence string) error {
	_, err := s.db.ExecContext(ctx, `INSERT INTO dna (user, sequence) VALUES (?,?)`, user, sequence)
	if err != nil {
		return fmt.Errorf("error writing to the repository")
	}

	return nil
}

// Select a user dna sequence from the database
func (s *SQLiteRepository) Select(ctx context.Context, user string) (sequence string, err error) {
	if err := s.db.QueryRowContext(ctx, `SELECT sequence FROM dna WHERE user = ?`, user).Scan(&sequence); err == sql.ErrNoRows {
		return "", ErrInvalidUser
	} else if err != nil {
		return "", fmt.Errorf("error reading repository: %v ", err)
	}
	return sequence, nil
}
