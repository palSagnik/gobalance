package main

import (
	"flag"
	"fmt"
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
	url *url.URL
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
	ServerList *ServerList
}

func NewOdin(conf *Config) *Odin {
	servers := make([]*Server, 0)
	for _, service := range conf.Services {
		for _, replica := range service.Replicas {
			serverUrl, err := url.Parse(replica)
			if err != nil {
				log.Fatal(err)
			}

			serverProxy := httputil.NewSingleHostReverseProxy(serverUrl)
			servers = append(servers, &Server{
				url: serverUrl,
				proxy: serverProxy,
			})
		}
	}

	return &Odin{
		Config: conf,
		ServerList: &ServerList{
			Servers: servers,
			currentServer: uint32(0),
		},
	}
}

func (o *Odin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Infof("received new request: url = '%s'", req.Host)
	next := o.ServerList.Next()
	o.ServerList.Servers[next].proxy.ServeHTTP(res, req)
}


func main() {

	flag.Parse()
	conf := &Config{
		Services: []Service {
			{
				Name: "Test1",
				Replicas: []string{"http://localhost:8081"},
			},
			{
				Name: "Test2",
				Replicas: []string{"http://localhost:8082"},
			},
			{
				Name: "Test3",
				Replicas: []string{"http://localhost:8083"},
			},
		},
	}
	odin := NewOdin(conf)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: odin,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
