package dna

import "testing"

func TestFlow(t *testing.T) {
	var (
		repo  = newMockRepo()
		user  = "vincent"
		token = "some_token"
		valid = newMockValidator(user, token)
		s     = NewDefaultService(repo, valid)
	)
}

type mockRepo struct {
	dna map[string]string
}

func newMockRepo() *mockRepo {
	return &mockRepo{
		dna: map[string]string{},
	}
}

type mockValidator struct {
	tokens map[string]string
}

func newMockValidator(usertokens ...string) *mockValidator {
	tokens := map[string]string{}
	for i := 0; i < len(usertokens); i += 2 {
		tokens[usertokens[i]] = usertokens[i+2]
	}
	return &mockValidator{
		tokens: tokens,
	}
}
