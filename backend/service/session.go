package service

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/blue-axes/tmpl/types"
)

type session struct {
	token     string
	userInfo  *types.UserInfo
	expiresAt time.Time
}

type SessionStore struct {
	ttl      time.Duration
	mu       sync.RWMutex
	sessions map[string]*session
}

func NewSessionStore(ttlHours int) *SessionStore {
	s := &SessionStore{
		ttl:      time.Duration(ttlHours) * time.Hour,
		sessions: make(map[string]*session),
	}
	go s.cleanupLoop()
	return s
}

func (s *SessionStore) Create(userInfo *types.UserInfo) (string, time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	token := generateToken()
	s.sessions[token] = &session{
		token:     token,
		userInfo:  userInfo,
		expiresAt: time.Now().Add(s.ttl),
	}
	return token, s.ttl
}

func (s *SessionStore) Validate(token string) *types.UserInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, ok := s.sessions[token]
	if !ok || time.Now().After(sess.expiresAt) {
		return nil
	}
	return sess.userInfo
}

func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

func (s *SessionStore) cleanupLoop() {
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		s.mu.Lock()
		now := time.Now()
		for token, sess := range s.sessions {
			if now.After(sess.expiresAt) {
				delete(s.sessions, token)
			}
		}
		s.mu.Unlock()
	}
}

func generateToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}