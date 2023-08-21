package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"dns/internal/api/domains"
	"dns/internal/dns"
)

func main() {
	d := domains.Controller{
		DNS: dns.NewService(),
	}

	r := chi.NewRouter()
	r.Get("/api/v1/resolvers/{resolver}/domains/{domain}", d.Show)
	r.HandleFunc("/resolvers/{resolver}/domains/{domain}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
	r.Handle("/*", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServe(":80", r))
}
