package config

import (
	"io"

	"github.com/palSagnik/gobalance/pkg/domain"
	"gopkg.in/yaml.v3"
)

func LoadConfig(r io.Reader) (*domain.Config, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	conf := domain.Config{}
	if err := yaml.Unmarshal(buf, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}