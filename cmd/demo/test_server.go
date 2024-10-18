package main

import (
	"flag"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

var (
	port = flag.Int("port", 8081, "Default port to start demo service on")
)

type DemoServer struct {}

func (d *DemoServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	body := []byte(fmt.Sprintf("all good from the server: %d\n", *port))
	
	res.WriteHeader(200)
	res.Write(body)
}

func main() {
	flag.Parse()

	addr := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(addr, &DemoServer{}); err != nil {
		log.Fatal(err)
	}
}