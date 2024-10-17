package main

import (
	"flag"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "port to start lb")
)

type Service struct {
	Name     string
	Replicas []string
}

// This is a representation of a configuration given to the loadbalancer
type Config struct {
	Services []Service

	// Name of the strategy used for load balancing
	Strategy string
}

// Server is an instance of a running server
type Server struct{
	url url.URL
	proxy *httputil.ReverseProxy
}

type ServerList struct{
	Servers []*Server

	// The list of servers are circulated through in a cyclic manner
	// next server is (currentServer + 1) * len(servers)
	currentServer uint32
}

func (sl *ServerList) Next() uint32 {
	nxt := atomic.AddUint32(&sl.currentServer, uint32(1))
	return nxt % uint32(len(sl.Servers))
}

type Odin struct {
	Config *Config
}

func (o *Odin) ServeHTTP(res http.ResponseWriter, req *http.Request) {

}
func main() {

	flag.Parse()
	conf := &Config{}
	odin := &Odin{
		Config: conf,
	}

	server := http.Server{
		Addr:    "",
		Handler: odin,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
