package dna

import (
	"database/sql"
	"errors"
	"fmt"
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
