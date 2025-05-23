package credence

import (
	"context"
	"fmt"
	"golang.org/x/oauth2"
	"sync"
)

// Config holds the configuration for a credential source
// such as a client ID/secret pair and token endpoint.
type Config struct {
	ClientID     string
	ClientSecret string
	TokenURL     string
	Scopes       []string
	Endpoint     oauth2.Endpoint // optional override
}

// TokenStore defines an interface for caching and retrieving tokens.
type TokenStore interface {
	Get(ctx context.Context, key string) (*oauth2.Token, error)
	Set(ctx context.Context, key string, token *oauth2.Token) error
}

// tokenEntry wraps the token and its sync lock.
type tokenEntry struct {
	token *oauth2.Token
	mu    sync.Mutex
}

// manager manages configs and token lifecycles.
type manager struct {
	configs map[string]Config
	store   TokenStore
	mu      sync.RWMutex
}

var global = &manager{
	configs: make(map[string]Config),
	store:   &memoryStore{cache: make(map[string]*oauth2.Token)},
}

// Register adds a new credential configuration under a key.
func Register(key string, cfg Config) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.configs[key] = cfg
}

// Token returns a valid access token for the given key.
func Token(ctx context.Context, key string) (string, error) {
	global.mu.RLock()
	cfg, ok := global.configs[key]
	global.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("%w", ErrConfigNotFound)
	}

	tok, err := global.store.Get(ctx, key)
	if err == nil && tok.Valid() {
		return tok.AccessToken, nil
	}

	oauthCfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     cfg.Endpoint,
		Scopes:       cfg.Scopes,
	}

	tokSrc := oauthCfg.TokenSource(ctx, nil)
	tok, err = tokSrc.Token()
	if err != nil {
		return "", err
	}

	_ = global.store.Set(ctx, key, tok)
	return tok.AccessToken, nil
}

// memoryStore is an in-memory TokenStore implementation.
type memoryStore struct {
	cache map[string]*oauth2.Token
	mu    sync.RWMutex
}

func (m *memoryStore) Get(ctx context.Context, key string) (*oauth2.Token, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tok, ok := m.cache[key]
	if !ok {
		return nil, ErrTokenNotFound
	}
	return tok, nil
}

func (m *memoryStore) Set(ctx context.Context, key string, token *oauth2.Token) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache[key] = token
	return nil
}

// Error types
var (
	ErrConfigNotFound = &Error{Msg: "credential config not found"}
	ErrTokenNotFound  = &Error{Msg: "token not found"}
)

type Error struct {
	Msg string
}

func (e *Error) Error() string {
	return e.Msg
}
