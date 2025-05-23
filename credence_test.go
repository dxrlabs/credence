package credence_test

import (
	"context"
	"errors"
	"github.com/dxrlabs/credence"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

type mockStore struct {
	tokens map[string]*oauth2.Token
}

func newMockStore() *mockStore {
	return &mockStore{
		tokens: make(map[string]*oauth2.Token),
	}
}

func (m *mockStore) Get(ctx context.Context, key string) (*oauth2.Token, error) {
	tok, ok := m.tokens[key]
	if !ok {
		return nil, credence.ErrTokenNotFound
	}
	return tok, nil
}

func (m *mockStore) Set(ctx context.Context, key string, token *oauth2.Token) error {
	m.tokens[key] = token
	return nil
}

func TestMemoryStore_SetGet(t *testing.T) {
	store := newMockStore()
	ctx := context.Background()
	tok := &oauth2.Token{
		AccessToken: "abc123",
		Expiry:      time.Now().Add(time.Hour),
	}

	err := store.Set(ctx, "service", tok)
	if err != nil {
		t.Fatalf("unexpected error setting token: %v", err)
	}

	got, err := store.Get(ctx, "service")
	if err != nil {
		t.Fatalf("unexpected error getting token: %v", err)
	}
	if got.AccessToken != "abc123" {
		t.Errorf("expected token 'abc123', got '%s'", got.AccessToken)
	}
}

func TestToken_MissingConfig(t *testing.T) {
	ctx := context.Background()
	_, err := credence.Token(ctx, "nonexistent")
	if !errors.Is(err, credence.ErrConfigNotFound) {
		t.Errorf("expected ErrConfigNotFound, got %v", err)
	}
}
