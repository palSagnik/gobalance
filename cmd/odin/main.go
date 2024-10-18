package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	models "github.com/palSagnik/gobalance/pkg/models"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "port to start lb")
)

type Odin struct {
	Config *models.Config
	ServerList *models.ServerList
}

func NewOdin(conf *models.Config) *Odin {
	servers := make([]*models.Server, 0)
	for _, service := range conf.Services {
		for _, replica := range service.Replicas {
			serverUrl, err := url.Parse(replica)
			if err != nil {
				log.Fatal(err)
			}

			serverProxy := httputil.NewSingleHostReverseProxy(serverUrl)
			servers = append(servers, &models.Server{
				Url: serverUrl,
				Proxy: serverProxy,
			})
		}
	}

	return &Odin{
		Config: conf,
		ServerList: &models.ServerList{
			Servers: servers,
			CurrentServer: uint32(0),
		},
	}
}

func (o *Odin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Infof("received new request: url = '%s'", req.Host)

	next := o.ServerList.Next()
	log.Infof("forwarding to server: '%d'", next)
	o.ServerList.Servers[next].Proxy.ServeHTTP(res, req)
}


func main() {
	flag.Parse()

	conf := &models.Config{
		Services: []models.Service {
			{
				Name: "Test",
				Replicas: []string{"http://localhost:8081", "http://localhost:8082"},
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
