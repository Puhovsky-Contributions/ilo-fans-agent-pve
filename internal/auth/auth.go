package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Store struct {
	path string
	mu   sync.RWMutex
	token string
}

func NewStore(path string) *Store {
	return &Store{path: path}
}

func (s *Store) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := os.ReadFile(s.path)
	if err != nil {
		return err
	}
	tok := strings.TrimSpace(string(data))
	if tok == "" {
		return fmt.Errorf("token file empty")
	}
	s.token = tok
	return nil
}

func (s *Store) Token() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.token
}

func (s *Store) Valid(bearer string) bool {
	if bearer == "" {
		return false
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(bearer, prefix) {
		return false
	}
	got := strings.TrimSpace(strings.TrimPrefix(bearer, prefix))
	return got != "" && got == s.Token()
}

func (s *Store) Init() (string, error) {
	tok, err := generateToken()
	if err != nil {
		return "", err
	}
	if err := writeTokenFile(s.path, tok); err != nil {
		return "", err
	}
	s.mu.Lock()
	s.token = tok
	s.mu.Unlock()
	return tok, nil
}

func (s *Store) Rotate() (string, error) {
	return s.Init()
}

func (s *Store) ShowMasked() string {
	t := s.Token()
	if len(t) <= 8 {
		return "****"
	}
	return t[:4] + "…" + t[len(t)-4:]
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func writeTokenFile(path, token string) error {
	if err := os.MkdirAll(dirOf(path), 0o750); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(token+"\n"), 0o600)
}

func dirOf(path string) string {
	i := strings.LastIndex(path, "/")
	if i <= 0 {
		return "."
	}
	return path[:i]
}
