package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := LoadConfig(strings.NewReader(`
strategy: Round-Robin
service:
  - name: test service
    replicas:
      - localhost:8081
      - localhost:8082`,))
	if err != nil {
		t.Error(err)
	}
}