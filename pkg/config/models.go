package config

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

type Service struct {
	Name     string   `yaml:"name"`
	Replicas []string `yaml:"replicas"`
}

// This is a representation of a configuration given to the loadbalancer
type Config struct {
	Services []Service `yaml:"service"`
	Strategy string    `yaml:"strategy"` // Name of the strategy used for load balancing
}

// Server is an instance of a running server
type Server struct {
	Url   *url.URL
	Proxy *httputil.ReverseProxy
}

func (s *Server) Forward(res http.ResponseWriter, req *http.Request) {
	s.Proxy.ServeHTTP(res, req)
}

type ServerList struct {
	Servers []*Server

	// The list of servers are circulated through in a cyclic manner
	// next server is (currentServer + 1) * len(servers)
	CurrentServer uint32
}

func (sl *ServerList) Next() uint32 {
	nxt := atomic.AddUint32(&sl.CurrentServer, uint32(1))
	return nxt % uint32(len(sl.Servers))
}
