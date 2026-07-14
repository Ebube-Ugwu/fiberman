package fiberman

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort        string
	NodeURL           string
	AuthToken         string
	TimeoutSeconds    int64
	PlaygroundBaseURL string
}

func LoadConfigFromEnv() Config {
	return Config{
		ServerPort:        envOrDefault("SERVER_PORT", "9020"),
		NodeURL:           envOrDefault("FIBER_NODE_URL", "http://127.0.0.1:8227"),
		AuthToken:         strings.TrimSpace(os.Getenv("FIBER_NODE_AUTH_TOKEN")),
		TimeoutSeconds:    intEnvOrDefault("FIBER_NODE_TIMEOUT_SECONDS", 30),
		PlaygroundBaseURL: envOrDefault("FIBER_PLAYGROUND_BASE_URL", "http://localhost:9020"),
	}
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func intEnvOrDefault(key string, fallback int64) int64 {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed < 1 {
		return fallback
	}
	return parsed
}
