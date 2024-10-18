package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	config "github.com/palSagnik/gobalance/pkg/config"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "port to start odin")
	configfile = flag.String("config-file", "", "configuration file to supply to odin")
)

type Odin struct {
	Config *config.Config
	ServerList *config.ServerList
}

func NewOdin(conf *config.Config) *Odin {
	servers := make([]*config.Server, 0)
	for _, service := range conf.Services {
		for _, replica := range service.Replicas {
			serverUrl, err := url.Parse(replica)
			if err != nil {
				log.Fatal(err)
			}

			serverProxy := httputil.NewSingleHostReverseProxy(serverUrl)
			servers = append(servers, &config.Server{
				Url: serverUrl,
				Proxy: serverProxy,
			})
		}
	}

	return &Odin{
		Config: conf,
		ServerList: &config.ServerList{
			Servers: servers,
			CurrentServer: uint32(0),
		},
	}
}

func (o *Odin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Infof("received new request: url = '%s'", req.Host)

	next := o.ServerList.Next()
	log.Infof("forwarding to server: '%d'", next)
	o.ServerList.Servers[next].Forward(res, req)
}


func main() {
	flag.Parse()

	// handling file
	file, err := os.Open(*configfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// handling configuration
	conf, err := config.LoadConfig(file)
	if err != nil {
		log.Fatal(err)
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
