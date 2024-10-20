package strategy

import (
	"sync/atomic"

	"github.com/palSagnik/gobalance/pkg/config"
)


type BalancingStrategy interface {
	Next([]*config.Server) (*config.Server, error)
}


type RoundRobin struct {
	CurrentServer uint32
}

func (rr *RoundRobin) Next(givenServers []*config.Server) (*config.Server, error) {
	nxt := atomic.AddUint32(&rr.CurrentServer, uint32(1))
	lenS := uint32(len(givenServers))

	return givenServers[nxt % lenS], nil
}