package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	config "github.com/palSagnik/gobalance/pkg/config"
	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8080, "port to start odin")
	configfile = flag.String("config-file", "", "configuration file to supply to odin")
)

type Odin struct {
	Config *config.Config
	ServerList map[string]*config.ServerList
}

func NewOdin(conf *config.Config) *Odin {
	serverMap := make(map[string]*config.ServerList, 0)
	for _, service := range conf.Services {
		servers := make([]*config.Server, 0)
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
		serverMap[service.Matcher] = &config.ServerList{
			Servers: servers,
			CurrentServer: uint32(0),
			Name: service.Name,
		}
	}

	return &Odin{
		Config: conf,
		ServerList: serverMap,
	}
}

// finds the first server list which matches the req path
// returns an error if no match found
func (o *Odin) findServiceList(reqPath string) (*config.ServerList, error) {

	log.Infof("trying to find a matcher for the request: '%s'", reqPath)
	for matcher, s := range o.ServerList {
		if strings.HasPrefix(reqPath, matcher) {
			log.Infof("found the service '%s' for the matching request", s.Name)
			return s, nil
		}
	}
	return nil, fmt.Errorf("did not find any matching service for url '%s'", reqPath)
}


func (o *Odin) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Infof("received new request: url = '%s'", req.Host)

	sl, err := o.findServiceList(req.URL.Path)
	if err != nil {
		log.Error(err)
		res.WriteHeader(http.StatusNotFound)
		return
	}
	next := sl.Next()
	log.Infof("forwarding to server: '%d'", next)
	sl.Servers[next].Forward(res, req)
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
