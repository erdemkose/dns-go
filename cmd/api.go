package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"dns/internal/api/domains"
	"dns/internal/dns"
)

func main() {
	var unixSocket string
	flag.StringVar(&unixSocket, "unix-socket", "dns-api.sock", "Unix socket")
	flag.Parse()

	d := domains.Controller{
		DNS: dns.NewService(),
	}

	r := chi.NewRouter()
	r.Get("/v1/resolvers/{resolver}/domains/{domain}", d.Show)

	l, err := net.Listen("unix", unixSocket)
	if err != nil {
		log.Fatal(err)
	}

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(l)

	s := http.Server{Handler: r}
	log.Fatal(s.Serve(l))
}
