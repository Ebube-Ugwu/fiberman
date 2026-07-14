package fiberman

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

const (
	sessionCookieName = "fiberman-go-session"
	maxHistorySize    = 20
)

type SessionStore struct {
	mu        sync.RWMutex
	histories map[string][]FiberCallResponse
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		histories: make(map[string][]FiberCallResponse),
	}
}

func (s *SessionStore) SessionID(w http.ResponseWriter, r *http.Request) string {
	if cookie, err := r.Cookie(sessionCookieName); err == nil && cookie.Value != "" {
		return cookie.Value
	}

	sessionID := randomID()
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	return sessionID
}

func (s *SessionStore) Append(sessionID string, entry FiberCallResponse) FiberCallResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	history := s.histories[sessionID]
	history = append([]FiberCallResponse{entry}, history...)
	if len(history) > maxHistorySize {
		history = history[:maxHistorySize]
	}
	s.histories[sessionID] = history
	return entry
}

func (s *SessionStore) History(sessionID string) []FiberCallResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := s.histories[sessionID]
	cloned := make([]FiberCallResponse, len(history))
	copy(cloned, history)
	return cloned
}

func (s *SessionStore) Clear(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.histories, sessionID)
}

func randomID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "session-fallback"
	}
	return hex.EncodeToString(bytes)
}
