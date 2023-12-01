package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"dns/internal/api/domains"
	"dns/internal/dns"
)

func main() {
	d := domains.Controller{
		DNS: dns.NewService(),
	}

	c := cors.New(cors.Options{
		AllowedOrigins: []string{
			"https://dnschecker.app/*",
			"https://dnschecker.pages.dev/*",
		},
		AllowedMethods: []string{
			"GET",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Language",
			"Content-Type",
		},
		ExposedHeaders: []string{
			"Content-Language",
			"Content-Length",
			"Content-Type",
			"Expires",
			"Last-Modified",
			"Pragma",
		},
		MaxAge: 300,
	})

	r := chi.NewRouter()
	r.Use(c.Handler)
	r.Get("/v1/resolvers/{resolver}/domains/{domain}", d.Show)

	network, address, ok := parseNetworkAndAddress()
	if !ok {
		log.Fatal("missing or invalid port and unix socket parameter")
	}

	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on %s://%s\n", network, address)

	defer func(l net.Listener) {
		err := l.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(l)

	s := http.Server{Handler: r}
	log.Fatal(s.Serve(l))
}

func parseNetworkAndAddress() (string, string, bool) {
	var port, unixSocket string
	flag.StringVar(&port, "port", "", "TCP Port")
	flag.StringVar(&unixSocket, "unix-socket", "", "Unix socket")
	flag.Parse()

	if port == "" {
		if p, found := os.LookupEnv("PORT"); found {
			port = p
		}
	}

	if port != "" {
		return "tcp", "0.0.0.0:" + port, true
	}

	if unixSocket != "" {
		return "unix", unixSocket, true
	}

	return "", "", false
}
