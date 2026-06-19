package main

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

// Session tracks a single PAM conversation
type Session struct {
	ID         string
	Username   string
	PromptCh   chan string
	ResponseCh chan string
	ResultCh   chan error
	CreatedAt  time.Time
}

// SessionStore manages all active sessions
type SessionStore struct {
	Sessions map[string]*Session
	Mu       sync.RWMutex
	Timeout  time.Duration
}

// NewSessionStore creates a new session store with the given timeout
func NewSessionStore(timeout time.Duration) *SessionStore {
	return &SessionStore{
		Sessions: make(map[string]*Session),
		Timeout:  timeout,
	}
}

// Create creates a new session and returns it
func (s *SessionStore) Create(username string) *Session {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	session := &Session{
		ID:         generateSessionID(),
		Username:   username,
		PromptCh:   make(chan string, 10),
		ResponseCh: make(chan string),
		ResultCh:   make(chan error),
		CreatedAt:  time.Now(),
	}
	s.Sessions[session.ID] = session
	return session
}

// Get retrieves a session by ID
func (s *SessionStore) Get(id string) *Session {
	s.Mu.RLock()
	defer s.Mu.RUnlock()
	return s.Sessions[id]
}

// Delete removes a session by ID
func (s *SessionStore) Delete(id string) {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	if session, exists := s.Sessions[id]; exists {
		close(session.PromptCh)
		close(session.ResponseCh)
		close(session.ResultCh)
		delete(s.Sessions, id)
	}
}

// Cleanup removes expired sessions (run in background goroutine)
func (s *SessionStore) Cleanup() {
	for {
		time.Sleep(s.Timeout / 2)
		s.Mu.Lock()
		now := time.Now()
		for id, session := range s.Sessions {
			if now.Sub(session.CreatedAt) > s.Timeout {
				close(session.PromptCh)
				close(session.ResponseCh)
				close(session.ResultCh)
				delete(s.Sessions, id)
			}
		}
		s.Mu.Unlock()
	}
}

// generateSessionID creates a random session ID
func generateSessionID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}
