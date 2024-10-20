package config

import (
	"github.com/palSagnik/gobalance/pkg/domain"
	"github.com/palSagnik/gobalance/pkg/strategy"
)

// This is a representation of a configuration given to the loadbalancer
type Config struct {
	Services []domain.Service `yaml:"services"`
	Strategy string    `yaml:"strategy"` // Name of the strategy used for load balancing
}

// This is the server list for a particular service
type ServerList struct {
	Servers  []*domain.Server                  	// Servers are the Replicas
	Name     string                     	// This is the name of the service in the configuration file
	Strategy strategy.BalancingStrategy 	// This is how the service is load balanced. It should never be empty and should default to RoundRobin
}


