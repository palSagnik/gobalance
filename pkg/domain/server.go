package domain

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Service struct {
	Name     string   `yaml:"name"`
	Replicas []string `yaml:"replicas"`
	Matcher  string   `yaml:"matcher"` // A prefix to select the service based on the path of the url
	Strategy string   `yaml:"strategy"` // Load balancing strategy used for this service
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
