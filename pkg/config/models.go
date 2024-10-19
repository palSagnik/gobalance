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
	Matcher  string   `yaml:"matcher"` // A prefix to select the service based on the path of the url
}

// This is a representation of a configuration given to the loadbalancer
type Config struct {
	Services []Service `yaml:"services"`
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

// This is the server list for a particular service
type ServerList struct {
	// Servers are the Replicas
	Servers []*Server

	// This is the name of the service in the configuration file
	Name    string

	CurrentServer uint32
}

func (sl *ServerList) Next() uint32 {
	nxt := atomic.AddUint32(&sl.CurrentServer, uint32(1))
	return nxt % uint32(len(sl.Servers))
}
