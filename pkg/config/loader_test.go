package config

import (
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	conf, err := LoadConfig(strings.NewReader(`
strategy: RoundRobin
services:
  - name: test service
    matcher: /api/v1
    replicas:
      - localhost:8081
      - localhost:8082`))
	if err != nil {
		t.Error(err)
	}

	if conf.Strategy != "RoundRobin" {
		t.Errorf("strategy expected to be `RoundRobin` found `%s` instead.", conf.Strategy)
	}
	if len(conf.Services) != 1 {
		t.Errorf("expected services count to be 1 got %d instead.", len(conf.Services))
	}
	if conf.Services[0].Name != "test service" {
		t.Errorf("expected service name to be equal to `test service` found %s instead.", conf.Services[0].Name)
	}
	if conf.Services[0].Matcher != "/api/v1" {
		t.Errorf("expected service name to be equal to `/api/v1` found %s instead.", conf.Services[0].Matcher)
	}
	if len(conf.Services[0].Replicas) != 2 {
		t.Errorf("expected service replicas to be equal to 2 found %d instead.", len(conf.Services[0].Replicas))
	}
}