package fiberman

import (
	"strings"
	"sync"
	"time"

	"github.com/fiberman/fiber-go-sdk/client"
)

type RuntimeSettings struct {
	mu                     sync.RWMutex
	nodeURL                string
	authToken              string
	timeoutSeconds         int64
	defaultInvoiceCurrency string
	playgroundBaseURL      string
}

func NewRuntimeSettings(config Config) *RuntimeSettings {
	return &RuntimeSettings{
		nodeURL:           config.NodeURL,
		authToken:         strings.TrimSpace(config.AuthToken),
		timeoutSeconds:    config.TimeoutSeconds,
		playgroundBaseURL: config.PlaygroundBaseURL,
	}
}

func (s *RuntimeSettings) Client() (*client.Client, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	timeout := time.Duration(s.timeoutSeconds) * time.Second
	return client.New(client.Config{
		BaseURL:        s.nodeURL,
		AuthToken:      s.authToken,
		ConnectTimeout: timeout,
		RequestTimeout: timeout,
	})
}

func (s *RuntimeSettings) Get() PlaygroundSettingsResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return PlaygroundSettingsResponse{
		NodeURL:                s.nodeURL,
		AuthToken:              s.authToken,
		TimeoutSeconds:         s.timeoutSeconds,
		DefaultInvoiceCurrency: s.defaultInvoiceCurrency,
		PlaygroundBaseURL:      s.playgroundBaseURL,
	}
}

func (s *RuntimeSettings) Update(request UpdatePlaygroundSettingsRequest) PlaygroundSettingsResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nodeURL = strings.TrimSpace(request.NodeURL)
	s.authToken = strings.TrimSpace(request.AuthToken)
	s.timeoutSeconds = request.TimeoutSeconds
	s.defaultInvoiceCurrency = strings.TrimSpace(request.DefaultInvoiceCurrency)

	return PlaygroundSettingsResponse{
		NodeURL:                s.nodeURL,
		AuthToken:              s.authToken,
		TimeoutSeconds:         s.timeoutSeconds,
		DefaultInvoiceCurrency: s.defaultInvoiceCurrency,
		PlaygroundBaseURL:      s.playgroundBaseURL,
	}
}

func (s *RuntimeSettings) NodeURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nodeURL
}

func (s *RuntimeSettings) AuthToken() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.authToken
}

func (s *RuntimeSettings) PlaygroundBaseURL() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.playgroundBaseURL
}
