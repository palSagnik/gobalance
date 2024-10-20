package strategy

import (
	"sync/atomic"

	"github.com/palSagnik/gobalance/pkg/domain"
	log "github.com/sirupsen/logrus"
)

// Load balancing strategies
const (
	RoundRobinStrategy = "RoundRobin"
	WeightedRoundRobinStrategy = "WeightedRoundRobin"
	UnknownStrategy = "Unknown"
)

type BalancingStrategy interface {
	Next([]*domain.Server) (*domain.Server, error)
}

var strategies map[string]func() BalancingStrategy
func init() {
	strategies = make(map[string]func() BalancingStrategy, 0)
	strategies[RoundRobinStrategy] = func() BalancingStrategy {
		return &RoundRobin{
			currentServer: uint32(0),
		}
	}
}

type RoundRobin struct {
	currentServer uint32
}

func (rr *RoundRobin) Next(givenServers []*domain.Server) (*domain.Server, error) {
	nxt := atomic.AddUint32(&rr.currentServer, uint32(1))
	lenS := uint32(len(givenServers))
	selectedServer := givenServers[nxt % lenS]
	log.Infof("Strategy selected server: '%s", selectedServer.Url.Host)
	return selectedServer, nil
}

// LoadStrategy will try to resolve the strategy based on the given name
// else it would default to RoundRobin
func LoadStrategy(name string) BalancingStrategy {
	st, ok := strategies[name]
	if !ok {
		return strategies[RoundRobinStrategy]()
	}

	return st()
}