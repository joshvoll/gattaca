package auth

import (
	"context"
	"testing"
)

func TestSQLiteFixture(t *testing.T) {
	r, err := NewSQLiteRepository("file:testdata/fixture.db")
	if err != nil {
		t.Fatal(err)
	}
	_, err = r.Auth(context.Background(), "bob", "bad password")
	if want, have := ErrBadAuth, err; want != have {
		t.Errorf("Auth with bad creds: want %v, have %v ", want, have)
	}

}
