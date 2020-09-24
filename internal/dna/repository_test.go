package dna

import (
	"context"
	"os"
	"testing"
)

func TestSQLiteFixture(t *testing.T) {
	r, err := NewSQLiteRepository("file:testdata/fixture.db")
	if err != nil {
		t.Fatalf("Error creating the repo: %v ", err)
	}
	_, err = r.Select(context.Background(), "invalid user")
	if want, have := ErrInvalidUser, err; want != have {
		t.Errorf("select with bad user: want %v, have %v ", want, have)
	}
	sequence, err := r.Select(context.Background(), "charlie")
	if want, have := error(nil), err; want != have {
		t.Errorf("select with bad user: want %v, have %v ", want, have)
	}
	if want, have := "aaaggactgcgcgccagttaagccctgttgtt", sequence; want != have {
		t.Errorf("select sequence: want %v, have %v ", want, have)
	}
}

func TestSQLIntegration(t *testing.T) {
	var (
		filevar  = "DNA_INTEGRATION_TEST_FILE"
		filename = os.Getenv(filevar)
	)
	if filename == "" {
		t.Skipf("skiping; set %s to run this test", filevar)
	}
	if _, err := os.Stat(filename); !os.IsNotExist(err) {
		t.Fatalf("%s: %v ", filename, err)
	}
	defer func() {
		if err := os.Remove(filename); err != nil {
			t.Errorf("removing %s:%v ", filename, err)
		}
	}()
	r, err := NewSQLiteRepository("file:" + filename)
	if err != nil {
		t.Fatalf("Error creating the repo: %v ", err)
	}
	var (
		user     = "vincent"
		sequence = "gattaca"
	)
	if want, have := error(nil), r.Insert(context.Background(), user, sequence); want != have {
		t.Fatalf("Error inserting user: want %v, have %v ", want, have)
	}
	selected, err := r.Select(context.Background(), user)
	if want, have := error(nil), err; want != have {
		t.Fatalf("select: want %v, have %v", want, have)
	}
	if want, have := sequence, selected; want != have {
		t.Fatalf("select: want %v, have %v ", want, have)
	}
}
