package domains

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"dns/internal/dns"
)

type Controller struct {
	DNS dns.Service
}

func (c *Controller) Show(rw http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")
	resolver := chi.URLParam(r, "resolver")

	records, err := c.DNS.Resolve(r.Context(), domain, resolver)
	if err != nil {
		log.Printf("DNS resolve failure: domain: %s, resolver:%s, err:%s\n", domain, resolver, err.Error())

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write([]byte(err.Error()))

		return
	}

	log.Printf("DNS resolve success: domain: %s, resolver:%s\n", domain, resolver)

	out, _ := json.Marshal(records)

	rw.Header().Set("Content-Type", "application/json")
	_, _ = rw.Write(out)
}
